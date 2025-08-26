package server

import (
	"errors"

	trace "go.opentelemetry.io/proto/otlp/trace/v1"
)

var (
	// ErrNoIngestorRegistered is raised when to do an operation
	// that requires a trace ingestor and it isn't registered
	ErrNoIngestorRegistered = errors.New("no trace ingestor registered")
)

// TraceIngestor is the function signature for processing trace data
type TraceIngestor func([]*trace.ResourceSpans) error

// Server defines the interface for the trace server
type Server interface {
	// Start the server
	Start(addr string) error

	// Stop the server
	Stop() error

	// RegisterTraceIngestor registers function for processing trace data
	RegisterTraceIngestor(TraceIngestor)
}
