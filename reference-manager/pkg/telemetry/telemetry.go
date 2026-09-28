package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"

	"github.com/MEPhI-C22-501/arch-info-system-monorepo/reference-manager/pkg/telemetry/log"
	"github.com/MEPhI-C22-501/arch-info-system-monorepo/reference-manager/pkg/telemetry/metrics"
)

type Telemetry struct {
	Log     *slog.Logger
	Metrics *metrics.Metrics
}

func NewTelemetry(cfg *Config) (*Telemetry, error) {
	res := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.DeploymentEnvironmentNameKey.String(cfg.Resource.Environment),
		semconv.ServiceNameKey.String(cfg.Resource.ServiceName),
	)

	log, err := log.NewLogger(cfg.Log.Options(log.WithResource(res))...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	metrics, err := metrics.NewMetrics(cfg.Metrics.Options(metrics.WithResource(res))...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return &Telemetry{
		Log:     log,
		Metrics: metrics,
	}, nil
}

func (t *Telemetry) Startup(ctx context.Context) error {
	if err := t.Metrics.Startup(ctx); err != nil {
		return fmt.Errorf("failed to startup metrics: %w", err)
	}

	return nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	if err := t.Metrics.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown metrics: %w", err)
	}

	return nil
}
