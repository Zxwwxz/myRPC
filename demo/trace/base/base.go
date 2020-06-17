package base

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"io"
)

func Init(service string) (opentracing.Tracer, io.Closer) {
	transport, err := zipkin.NewHTTPTransport(
		"http://60.205.218.189:9411/api/v1/spans",
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		fmt.Println("NewHTTPTransport err:",err)
		return nil,nil
	}
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
		ServiceName:service,
	}
	r := jaeger.NewRemoteReporter(transport)
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		fmt.Println("NewTracer err:",err)
		return nil,nil
	}
	return tracer, closer
}
