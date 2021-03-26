package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

//追踪模块
//自己搭建：搭建java，elasticsearch，zipkin，jaeger
//docker：docker run -d --name jaeger（测试环境）
//基础概念：
//traceId：整个调用过程的唯一id
//spanId：当前调用步骤的id
var (
	tracer opentracing.Tracer
	closer io.Closer
)

//初始化追踪
func InitTrace(serviceName, addr, sampleType string, sampleRate float64) (err error) {
	//获取配置，提供上报地址
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
	//创建追踪对象
	tracer, closer, err = cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return err
	}
	//保存追踪对象到全局中
	opentracing.SetGlobalTracer(tracer)
	return nil
}

func GetTracer()(opentracing.Tracer)  {
	return tracer
}

func Stop()(error)  {
	return closer.Close()
}

