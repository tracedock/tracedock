package server

import (
	"context"
	"net/http"
	"regexp"
	"time"
)

// HTTPServer implements Server interface for the HTTP protocol
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

	if contentType != "application/x-protobuf" && contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if match, _ := regexp.MatchString("^/v1/traces(/)?$", r.URL.Path); !match {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
