package meta

import (
	"context"
)

type ServerMeta struct {
	Env           string  //环境
	IDC           string  //机房
	Cluster       string  //集群
	ServerIP      string  //服务端ip
	ClientIP      string  //客户端ip
	ServiceName   string  //服务名
	ServiceMethod string  //服务方法
}

type ServerMetaContextKey struct{}

var s = ServerMetaContextKey{}

func GetServerMeta(ctx context.Context) *ServerMeta {
	meta,ok := ctx.Value(s).(*ServerMeta)
	if ok != true{
		meta = &ServerMeta{}
	}
	return meta
}

func SetServerMeta(ctx context.Context,meta *ServerMeta) context.Context{
	return context.WithValue(ctx, s, meta)
}