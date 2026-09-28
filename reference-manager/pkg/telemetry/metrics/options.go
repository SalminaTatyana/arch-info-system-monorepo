package metrics

import (
	"errors"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

type options struct {
	exporterOptions []otlpmetricgrpc.Option
	readerOptions   []sdkmetric.PeriodicReaderOption
	resource        *resource.Resource
}

func defaultOptions() *options {
	return &options{
		exporterOptions: make([]otlpmetricgrpc.Option, 0),
		readerOptions:   make([]sdkmetric.PeriodicReaderOption, 0),
		resource:        resource.Default(),
	}
}

type Option func(*options) error

func WithOtelCollectorEndpoint(v string) Option {
	return func(o *options) error {
		if v == "" {
			return errors.New("otel collector endpoint must not be empty")
		}

		o.exporterOptions = append(o.exporterOptions, otlpmetricgrpc.WithEndpoint(v))
		return nil
	}
}

func WithPushInterval(v time.Duration) Option {
	return func(o *options) error {
		if v <= 0 {
			return errors.New("push interval must be greater than zero")
		}

		o.readerOptions = append(o.readerOptions, sdkmetric.WithInterval(v))
		return nil
	}
}

func WithPushTimeout(v time.Duration) Option {
	return func(o *options) error {
		if v <= 0 {
			return errors.New("push timeout must be greater than zero")
		}

		o.readerOptions = append(o.readerOptions, sdkmetric.WithTimeout(v))
		return nil
	}
}

func WithInsecure(v bool) Option {
	return func(o *options) error {
		if v {
			o.exporterOptions = append(o.exporterOptions, otlpmetricgrpc.WithInsecure())
		}

		return nil
	}
}

func WithResource(res *resource.Resource) Option {
	return func(o *options) error {
		if res == nil {
			return errors.New("resource must not be nil")
		}

		o.resource = res
		return nil
	}
}
