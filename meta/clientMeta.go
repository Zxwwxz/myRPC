package meta

import (
	"context"
	"google.golang.org/grpc"
	registryBase "myRPC/registry/base"
)

type ClientMeta struct {
	//调用方IDC
	CallerIDC string
	//调用方集群
	CallerCluster string
	//调用方名字
	CallerName string
	//服务提供方IDC
	ServiceIDC string
	//服务提供方集群
	ServiceCluster string
	//服务提供方
	ServiceName string
	//服务方法
	ServiceMethod string
	//TraceID
	TraceID string
	//环境
	Env string
	//当前节点
	CurNode *registryBase.Node
	//历史选择节点
	RemainNodes map[int]*registryBase.Node
	//服务提供方的节点列表
	AllNodes map[int]*registryBase.Node
	//当前请求使用的连接
	Conn *grpc.ClientConn
}

type ClientMetaContextKey struct{}

func GetClientMeta(ctx context.Context) *ClientMeta {
	meta, ok := ctx.Value(ClientMetaContextKey{}).(*ClientMeta)
	if !ok {
		meta = &ClientMeta{}
	}
	return meta
}

func InitClientMeta(ctx context.Context, serviceName, serviceMethod, callerName string) context.Context {
	meta := &ClientMeta{
		ServiceMethod:   serviceMethod,
		ServiceName:     serviceName,
		CallerName:      callerName,
	}
	return context.WithValue(ctx, ClientMetaContextKey{}, meta)
}

