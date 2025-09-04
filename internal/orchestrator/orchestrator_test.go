package orchestrator

import (
	"testing"

	trace "go.opentelemetry.io/proto/otlp/trace/v1"

	"github.com/stretchr/testify/assert"
	"github.com/tracedock/tracedock/internal/config"
)

func Test_Ingestor_IngestTrace(t *testing.T) {
	var ingestor = NewIngestor(config.NewConfig())

	t.Run("should handle nil ResourceSpans", func(t *testing.T) {
		var rs *trace.ResourceSpans

		assert.NoError(t, ingestor.IngestTrace(rs))
	})

	t.Run("should handle nil ScopeSpans", func(t *testing.T) {
		var rs = &trace.ResourceSpans{
			ScopeSpans: nil,
		}

		assert.NoError(t, ingestor.IngestTrace(rs))
	})
}
