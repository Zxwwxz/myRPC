package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"myRPC/demo/trace/base"
	"net/http"
)

func main() {
	tracer, closer := base.Init("service_name")
	defer closer.Close()
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("service_span", ext.RPCServerOption(spanCtx))
		defer span.Finish()
		r.FormValue("service_resp_str")
	})
	_ = http.ListenAndServe(":8081", nil)
}

