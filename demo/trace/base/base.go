package base

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)
//搭建java，elasticsearch，zipkin，jaeger
func Init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort:"47.92.212.70:6831",
		},
		ServiceName:service,
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		fmt.Println("NewTracer err:",err)
		return nil,nil
	}
	return tracer, closer
}
