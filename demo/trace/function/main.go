package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"myRPC/demo/trace/base"
)

func main() {
	tracer, closer := base.Init("serName")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("spanName")
	span.SetTag("main_tag_key", "main_tag_value")
	span.LogFields(log.String("main_fields_key", "main_fields_value"))
	span.LogKV("main_kv_key", "main_kv_value")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	formatString(ctx)
	printHello(ctx)
	span.Finish()
}

func formatString(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()
	span.SetTag("formatString_tag_key", "formatString_tag_value")
	span.LogFields(log.String("formatString_fields_key", "formatString_fields_value"))
	span.LogKV("formatString_kv_key", "formatString_kv_value")
}

func printHello(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()
	span.SetTag("printHello_tag_key", "printHello_tag_value")
	span.LogFields(log.String("printHello_fields_key", "printHello_fields_value"))
	span.LogKV("printHello_kv_key", "printHello_kv_value")
}
