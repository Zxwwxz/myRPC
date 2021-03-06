package mwTrace

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
	logBase "myRPC/log/base"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	"myRPC/trace"
	"strings"
)

// 追踪中间件
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

func TraceServiceMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//上下文取出md
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
			defer serverSpan.Finish()
			serverSpan.SetTag("trace_id", trace.GetTraceId(ctx))
			logBase.Debug("TraceServiceMiddleware,serverSpan=%v",serverSpan)
			ctx = opentracing.ContextWithSpan(ctx, serverSpan)
			resp, err = next(ctx, req)
			return
		}
	}
}

func TraceClientMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//获取全局追踪器
			tracer := opentracing.GlobalTracer()
			//父级上下文
			var parentSpanCtx opentracing.SpanContext
			//父级span
			if parent := opentracing.SpanFromContext(ctx); parent != nil {
				parentSpanCtx = parent.Context()
			}
			opts := []opentracing.StartSpanOption{
				opentracing.ChildOf(parentSpanCtx),
				ext.SpanKindRPCClient,
				opentracing.Tag{Key: string(ext.Component), Value: "client_trace"},
				opentracing.Tag{Key: "trace_id", Value: trace.GetTraceId(ctx)},
			}
			clientMeta := meta.GetClientMeta(ctx)
			//子span
			clientSpan := tracer.StartSpan(clientMeta.ClientName, opts...)
			defer clientSpan.Finish()
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.Pairs()
			}
			//span注入到http头，传递的是map
			if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
				logBase.Warn("TraceClientMiddleware,tracer.Inject,err=%v",err)
			}
			ctx = metadata.NewOutgoingContext(ctx, md)
			ctx = metadata.AppendToOutgoingContext(ctx, "trace_id", trace.GetTraceId(ctx))
			//子span存入上下文
			ctx = opentracing.ContextWithSpan(ctx, clientSpan)
			resp, err = next(ctx, req)
			if err != nil {
				ext.Error.Set(clientSpan, true)
				clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
			}
			logBase.Debug("TraceClientMiddleware,clientSpan=%v",clientSpan)
			return
		}
	}
}

//客户端生成追踪id
func TraceIdClientMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			var traceId string
			//从上下文拿到md
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				vals, ok := md["trace_id"]
				if ok && len(vals) > 0 {
					traceId = vals[0]
				}
			}
			//没有追踪id生成一个
			if len(traceId) == 0 {
				traceId = trace.GenTraceId()
			}
			logBase.Debug("TraceIdClientMiddleware,traceId=%s",traceId)
			//保存追踪id到上下文中
			ctx = trace.WithTraceId(ctx, traceId)
			resp, err = next(ctx, req)
			return
		}
	}
}