package orchestrator

import (
	"fmt"

	trace "go.opentelemetry.io/proto/otlp/trace/v1"

	"github.com/tracedock/tracedock/internal/config"
	"github.com/tracedock/tracedock/internal/logger"
)

type Ingestor struct {
	Config *config.Config
}

func NewIngestor(config *config.Config) *Ingestor {
	return &Ingestor{config}
}

func (i *Ingestor) IngestTrace(rs *trace.ResourceSpans) error {
	if rs == nil {
		return nil
	}

	totalSpans := 0
	for _, ss := range rs.ScopeSpans {
		totalSpans += len(ss.Spans)
	}

	logger.Debug(fmt.Sprintf("trace with %d spans ingested", totalSpans))

	return nil
}
