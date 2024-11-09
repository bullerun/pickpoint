package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

func MustLoad(serviceName string) (func(), error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			//LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger), config.Metrics(prometheus.New()))
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return func() { closer.Close() }, nil
}
