package orchestrator

import (
	"fmt"

	trace "go.opentelemetry.io/proto/otlp/trace/v1"

	"github.com/tracedock/tracedock/internal/config"
	"github.com/tracedock/tracedock/internal/logger"
	"github.com/tracedock/tracedock/internal/storage"
)

type Ingestor struct {
	Config *config.Config
	Queue  *storage.Queue
}

func NewIngestor(config *config.Config, queue *storage.Queue) *Ingestor {
	return &Ingestor{config, queue}
}

func (i *Ingestor) IngestTrace(rs *trace.ResourceSpans) error {
	if rs == nil {
		return nil
	}

	logger.Info(fmt.Sprintf("enqueuing resource span with %#v scope spans", len(rs.ScopeSpans)))
	return i.Queue.Enqueue(rs)
}
