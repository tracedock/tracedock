package orchestrator

import (
	"errors"
	"math/rand"
	"testing"

	trace "go.opentelemetry.io/proto/otlp/trace/v1"

	"github.com/stretchr/testify/assert"
	"github.com/tracedock/tracedock/internal/config"
	"github.com/tracedock/tracedock/internal/storage"
)

func Test_Ingestor_IngestTrace(t *testing.T) {
	var ingestor = NewIngestor(config.NewConfig(), storage.NewQueue())

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

	t.Run("should return error when enqueue returns error", func(t *testing.T) {
		var queue = NewMockQueue(t)
		var ingestor = NewIngestor(config.NewConfig(), storage.NewQueue())

		var item = rand.Int()
		var err = errors.New("")

		queue.EXPECT().Enqueue(item).Return(err)

		assert.Equal(t, err, ingestor.Queue.Enqueue(item))
	})
}
