package mwTrace

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
	logBase "myRPC/log/base"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	"strings"
)

type metadataTextMap metadata.MD

func (m metadataTextMap) Set(key, val string) {
	encodedKey, encodedVal := encodeKeyValue(key, val)
	m[encodedKey] = []string{encodedVal}
}

func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, vv := range m {
		for _, v := range vv {
			if decodedKey, decodedVal, err := metadata.DecodeKeyValue(k, v); err == nil {
				if err = callback(decodedKey, decodedVal); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed decoding opentracing from gRPC metadata: %v", err)
			}
		}
	}
	return nil
}

func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, "-bin") {
		val := base64.StdEncoding.EncodeToString([]byte(v))
		v = string(val)
	}
	return k, v
}


func TraceMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				md = metadata.Pairs()
			}
			tracer := opentracing.GlobalTracer()
			parentSpanContext, _ := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
			serverMeta := meta.GetServerMeta(ctx)
			serverSpan := tracer.StartSpan(
				serverMeta.ServiceMethod,
				ext.RPCServerOption(parentSpanContext),
				ext.SpanKindRPCServer,
			)
			serverSpan.SetTag("trace_id", logBase.GetTraceId(ctx))
			ctx = opentracing.ContextWithSpan(ctx, serverSpan)
			resp, err = next(ctx, req)
			serverSpan.Finish()
			return
		}
	}
}