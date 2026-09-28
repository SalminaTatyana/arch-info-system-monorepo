package metrics

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var ErrNotStarted = errors.New("metrics provider is not started")

type Metrics struct {
	options *options

	mu       sync.RWMutex
	provider *sdkmetric.MeterProvider
}

func NewMetrics(opts ...Option) (*Metrics, error) {
	options := defaultOptions()
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to initialize metrics options: %w", err)
		}
	}

	return &Metrics{
		options: options,
	}, nil
}

func (m *Metrics) Startup(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.provider != nil {
		return nil
	}

	exporter, err := otlpmetricgrpc.New(ctx, m.options.exporterOptions...)
	if err != nil {
		return fmt.Errorf("failed to initialize metrics exporter: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter, m.options.readerOptions...)
	m.provider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetric.WithResource(m.options.resource))

	return nil
}

func (m *Metrics) Meter(name string, opts ...otelmetric.MeterOption) (otelmetric.Meter, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.provider == nil {
		return nil, ErrNotStarted
	}

	return m.provider.Meter(name, opts...), nil
}

func (m *Metrics) ForceFlush(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.provider == nil {
		return ErrNotStarted
	}

	if err := m.provider.ForceFlush(ctx); err != nil {
		return fmt.Errorf("failed to flush metrics: %w", err)
	}

	return nil
}

func (m *Metrics) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.provider == nil {
		return nil
	}

	if err := m.provider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown metrics provider: %w", err)
	}
	m.provider = nil

	return nil
}
