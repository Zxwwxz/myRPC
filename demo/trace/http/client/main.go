package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"myRPC/demo/trace/base"
	"net/http"
	"net/url"
)

func main()  {
	tracer, closer := base.Init("client_name")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("client_span_name")
	span.SetTag("client_tag_key", "client_tag_value")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	formatString(ctx)
}

func formatString(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()
	v := url.Values{}
	v.Set("client_req_key", "client_req_value")
	span.SetTag("req_tag_key", "req_tag_value")
	url := "http://localhost:8081/hello?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("NewRequest err:",err)
		return
	}
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	resp, err := httpDo(req)
	if err != nil {
		fmt.Println("httpDo err:",err)
		return
	}
	helloStr := string(resp)
	span.LogFields(log.String("resp_fields", helloStr))
	span.LogKV("resp_kv", helloStr)
}

