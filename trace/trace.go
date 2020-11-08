package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
)

func InitTrace(serviceName, addr, sampleType string, sampleRate float64) (err error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: sampleRate,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort:addr,
		},
		ServiceName:serviceName,
	}
	tracer, closer, err = cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	return nil
}

func GetTracer()(opentracing.Tracer)  {
	return tracer
}

func Stop()(error)  {
	return closer.Close()
}

