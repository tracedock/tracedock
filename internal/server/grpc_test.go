package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tracecollectorv1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	trace "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestNewGRPCServer(t *testing.T) {
	server := NewGRPCServer()

	assert.NotNil(t, server)
	assert.NotNil(t, server.server)
}

func TestGRPCServer_Start(t *testing.T) {
	t.Run("should return error when no ingestor is registered", func(t *testing.T) {
		server := NewGRPCServer()
		err := server.Start(":8080")

		assert.Error(t, err)
		assert.Equal(t, ErrNoIngestorRegistered, err)
	})

	t.Run("should start server successfully when ingestor is registered", func(t *testing.T) {
		var addr = "0.0.0.0:8081"
		var done = make(chan error)

		var ingestor = func([]*trace.ResourceSpans) error { return nil }

		server := NewGRPCServer()
		server.RegisterTraceIngestor(ingestor)

		go func() {
			err := server.Start(addr)
			done <- err
		}()

		time.Sleep(100 * time.Millisecond)

		select {
		case err := <-done:
			assert.NoError(t, err)
		case <-time.After(100 * time.Millisecond):
			server.Stop()
		}
	})
}

func TestGRPCServer_Export(t *testing.T) {
	t.Run("should return error when no ingestor is registered", func(t *testing.T) {
		server := NewGRPCServer()

		req := &tracecollectorv1.ExportTraceServiceRequest{
			ResourceSpans: []*trace.ResourceSpans{
				{},
			},
		}

		_, err := server.Export(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("should process traces successfully", func(t *testing.T) {
		var processedTraces []*trace.ResourceSpans

		ingestor := func(resource []*trace.ResourceSpans) error {
			processedTraces = append(processedTraces, resource...)
			return nil
		}

		server := NewGRPCServer()
		server.RegisterTraceIngestor(ingestor)

		req := &tracecollectorv1.ExportTraceServiceRequest{
			ResourceSpans: []*trace.ResourceSpans{
				{},
				{},
			},
		}

		resp, err := server.Export(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, processedTraces, 2)
	})

	t.Run("should return error when ingestor fails", func(t *testing.T) {
		ingestor := func(resource []*trace.ResourceSpans) error {
			return assert.AnError
		}

		server := NewGRPCServer()
		server.RegisterTraceIngestor(ingestor)

		req := &tracecollectorv1.ExportTraceServiceRequest{
			ResourceSpans: []*trace.ResourceSpans{
				{},
			},
		}

		_, err := server.Export(context.Background(), req)
		assert.Error(t, err)
	})
}
