package server

import (
	"net/http/httptest"
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

		var ingestor = func(traces *trace.ResourceSpans) error {
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

		var ingestor = func(*trace.ResourceSpans) error {
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

func Test_HTTPServer_HandleRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		urlPath        string
		expectedStatus int
	}{
		{
			name:           "should return 415 for invalid content-type",
			method:         "POST",
			contentType:    "text/plain",
			urlPath:        "/v1/traces",
			expectedStatus: 415,
		},
		{
			name:           "should return 405 for invalid method",
			method:         "GET",
			contentType:    "application/json",
			urlPath:        "/v1/traces",
			expectedStatus: 405,
		},
		{
			name:           "should return 404 for inexistent path",
			method:         "POST",
			contentType:    "application/json",
			urlPath:        "/v1/invalid",
			expectedStatus: 404,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := NewHTTPServer()

			server.RegisterTraceIngestor(func(*trace.ResourceSpans) error {
				return nil
			})

			req := httptest.NewRequest(tc.method, tc.urlPath, nil)
			req.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			server.HandleRequest(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			assert.Equal(t, tc.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
