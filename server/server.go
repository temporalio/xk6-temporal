package server

import (
	"github.com/DataDog/temporalite"
	"go.temporal.io/server/common/metrics"
)

func NewTemporaliteServer(metricsHandler metrics.MetricsHandler) (*temporalite.Server, error) {
	s, err := temporalite.NewServer(
		temporalite.WithPersistenceDisabled(),
		temporalite.WithCustomMetricsHandler(metricsHandler),
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}
