package server

import (
	"context"
	"net/http"
	"time"
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

func (s *HTTPServer) handleRequest(w http.ResponseWriter, r *http.Request) {}
