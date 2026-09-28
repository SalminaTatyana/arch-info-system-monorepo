package metrics

import "time"

type Config struct {
	OtelCollectorEndpoint string        `yaml:"otel_collector_endpoint"`
	PushInterval          time.Duration `yaml:"push_interval"`
	PushTimeout           time.Duration `yaml:"push_timeout"`
}

func (c Config) Options(external ...Option) []Option {
	opts := []Option{
		WithOtelCollectorEndpoint(c.OtelCollectorEndpoint),
		WithPushInterval(c.PushInterval),
		WithPushTimeout(c.PushTimeout),
		WithInsecure(true),
	}

	return append(opts, external...)
}
