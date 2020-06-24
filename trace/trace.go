package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

func Init(serviceName, reportAddr, sampleType string, rate float64) (err error) {
	transport, err := zipkin.NewHTTPTransport(
		reportAddr,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		return err
	}
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: rate,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
		ServiceName:serviceName,
	}
	r := jaeger.NewRemoteReporter(transport)
	tracer, _, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	return nil
}

