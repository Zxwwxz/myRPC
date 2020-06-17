package main

import (
	"github.com/opentracing/opentracing-go/log"
	"myRPC/demo/trace/base"
)

func main() {
	tracer, closer := base.Init("serName")
	defer closer.Close()
	span := tracer.StartSpan("spanName")
	span.SetTag("tag_key", "tag_value")
	span.LogFields(log.String("log_fields_key", "log_fields_value"))
	span.LogKV("log_kv_key", "log_kv_value")
	span.Finish()
}
