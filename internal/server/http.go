package server

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tracedock/tracedock/internal/logger"
	prototrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	trace "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// HTTPServer implements Server interface for the HTTP protocol
//
// Notice: It isn't implementing 100% of the OpenTelemetry HTTP specification
// regarding to the responses bodies, instead, for now it only respond with
// correct status code without any response bodies.
//
// This behaviour was tested with Ruby and Python OpenTelemetry SDKs and worked
// with no issues.
//
// For more details: https://opentelemetry.io/docs/specs/otlp/#otlphttp-response
type HTTPServer struct {
	httpServer    *http.Server
	traceIngestor TraceIngestor
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

// Start the HTTP server
func (s *HTTPServer) Start(addr string) error {
	logger.Info(fmt.Sprintf("starting HTTP server at %s", addr))

	if s.traceIngestor == nil {
		return ErrNoIngestorRegistered
	}

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(s.HandleRequest),
	}

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop the HTTP server
func (s *HTTPServer) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	background := context.Background()
	ctx, cancel := context.WithTimeout(background, 5*time.Second)

	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// RegisterTraceIngestor registers a TraceIngestor function that will process all the
// incoming trace data
func (s *HTTPServer) RegisterTraceIngestor(ingestor TraceIngestor) {
	s.traceIngestor = ingestor
}

// HandleRequest handles incoming HTTP requests in order to get trace ingested
func (s *HTTPServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var contentType = r.Header.Get("Content-Type")

	w.Header().Set("Content-Type", contentType)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if match, _ := regexp.MatchString("^/v1/traces(/)?$", r.URL.Path); !match {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch contentType {
	case "application/json":
		s.HandleRequestWithJSON(w, r)
		return

	case "application/x-protobuf":
		s.HandleRequestWithProtobuf(w, r)
		return

	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

// HandleRequestWithJSON handles requests with content-type equals application/json
func (s *HTTPServer) HandleRequestWithJSON(w http.ResponseWriter, r *http.Request) {
	var resourceSpans trace.ResourceSpans

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := protojson.Unmarshal(reqBody, &resourceSpans); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.traceIngestor(&resourceSpans); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// HandleRequestWithJSON handles requests with content-type equals application/protobuf
func (s *HTTPServer) HandleRequestWithProtobuf(w http.ResponseWriter, r *http.Request) {
	var reqBodyReader = r.Body
	var resourceSpans prototrace.ExportTraceServiceRequest
	var contentEncoding = r.Header.Get("Content-Encoding")

	if strings.Contains(contentEncoding, "gzip") {
		gzr, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer gzr.Close()
		reqBodyReader = gzr
	}

	reqBody, err := io.ReadAll(reqBodyReader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := proto.Unmarshal(reqBody, &resourceSpans); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, rs := range resourceSpans.ResourceSpans {
		if err := s.traceIngestor(rs); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
