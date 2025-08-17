package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	trace "go.opentelemetry.io/proto/otlp/trace/v1"
)

func Test_HTTPServer_Start(t *testing.T) {
	t.Run("should return error when no ingestor is registered", func(t *testing.T) {
		err := NewHTTPServer().Start(":8080")

		assert.Error(t, err)
		assert.Equal(t, ErrNoIngestorRegistered, err)
	})

	t.Run("should return error when server fails to start", func(t *testing.T) {
		var addr = "0.0.0.0:8080"

		var ingestor = func(trace *trace.ResourceSpans) error {
			return nil
		}

		server := NewHTTPServer()
		server.RegisterTraceIngestor(ingestor)

		go func() { server.Start(addr) }()

		time.Sleep(100 * time.Millisecond)

		assert.Error(t, server.Start(addr))
		assert.NoError(t, server.Stop())
	})

	t.Run("should start server successfully when ingestor is registered", func(t *testing.T) {
		var addr = "0.0.0.0:8080"

		var done = make(chan error)

		var ingestor = func(trace *trace.ResourceSpans) error {
			return nil
		}

		server := NewHTTPServer()
		server.RegisterTraceIngestor(ingestor)

		go func() {
			time.Sleep(100 * time.Millisecond)

			err := server.Start(addr)
			done <- err
		}()

		select {
		case err := <-done:
			// The fact we are getting here means the server
			// failed to start so it should fail anyway
			assert.NoError(t, err)

		case <-time.After(100 * time.Millisecond):
			server.Stop()
			assert.Equal(t, addr, server.httpServer.Addr)
		}
	})
}
