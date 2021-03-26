package meta

import (
	"context"
)
//服务端中间件参数
type ServerMeta struct {
	Env           string  //环境
	IDC           string  //机房
	ServeiceIP    string  //服务端ip
	ServiceName   string  //服务名
	ServiceMethod string  //服务方法
}

type ServerMetaContextKey struct{}

var s = ServerMetaContextKey{}

func GetServerMeta(ctx context.Context) *ServerMeta {
	meta,ok := ctx.Value(s).(*ServerMeta)
	if !ok {
		meta = &ServerMeta{}
	}
	return meta
}

func SetServerMeta(ctx context.Context,meta *ServerMeta) context.Context{
	return context.WithValue(ctx, s, meta)
}