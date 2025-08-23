package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"errors"
)

var (
	ErrHTTPUnsupportedMethod  = errors.New("invalid HTTP method")
	ErrHTTPNotFound           = errors.New("not found")
	ErrHTTPUnsupported        = errors.New("unsupported content type")
	ErrHTTPInvalidContentType = errors.New("invalid content type")
)

type HTTPServer struct {
	httpServer    *http.Server
	traceIngestor TraceIngestor
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start(addr string) error {
	if s.traceIngestor == nil {
		return ErrNoIngestorRegistered
	}

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(s.handleRequest),
	}

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HTTPServer) Stop() error {
	background := context.Background()
	ctx, cancel := context.WithTimeout(background, 5*time.Second)

	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func (s *HTTPServer) RegisterTraceIngestor(ingestor TraceIngestor) {
	s.traceIngestor = ingestor
}

func (s *HTTPServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	contentType, err := s.validateRequest(r)

	if err != nil {
		s.errToHTTPResponse(err, w, r)
		return
	}

	if err := s.processRequest(contentType, r.Body); err != nil {
		s.errToHTTPResponse(err, w, r)
		return
	}
}

func (s *HTTPServer) validateRequest(r *http.Request) (contentType string, _ error) {
	contentType = r.Header.Get("Content-Type")

	if contentType != "application/x-protobuf" && contentType != "application/json" {
		return contentType, ErrHTTPInvalidContentType
	}

	if r.Method != http.MethodPost {
		return contentType, ErrHTTPUnsupportedMethod
	}

	if match, _ := regexp.MatchString("^/v1/traces(/)?$", r.URL.Path); !match {
		return contentType, ErrHTTPNotFound
	}

	return contentType, nil
}

func (s *HTTPServer) processRequest(_ string, _ io.ReadCloser) error {
	return nil
}

func (s *HTTPServer) errToHTTPResponse(err error, w http.ResponseWriter, r *http.Request) {
	switch {
	case errors.Is(err, ErrHTTPUnsupportedMethod):
		http.Error(w, fmt.Sprintf("unsupported HTTP method: %v", r.Method), http.StatusMethodNotAllowed)
	case errors.Is(err, ErrHTTPNotFound):
		http.NotFound(w, r)
	case errors.Is(err, ErrHTTPInvalidContentType):
		http.Error(w, fmt.Sprintf("invalid content type: %v", r.Header.Get("Content-Type")), http.StatusUnsupportedMediaType)
	default:
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
	}
}
