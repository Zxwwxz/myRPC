package meta

import (
	"context"
	"google.golang.org/grpc"
	"myRPC/registry/register"
)

const (
	//不指定哪个节点，由均衡算法算出
	caller_type_not = 1
	//指定服务类型，调用所有节点
	caller_type_type = 2
	//指定服务类型和服务id，调用指定节点
	caller_type_type_id = 3
)

type ClientMeta struct {
	//服务提供方
	ServiceName string
	//服务方法
	ServiceMethod string

	//调用类型
	CallerType int
	//当前节点
	CurNode *register.Node
	//历史选择节点
	RemainNodes []*register.Node
	//服务提供方的节点列表
	AllNodes []*register.Node
	//调用失败最大重连次数
	MaxReconnectNum int
	//负载均衡关键字
	BalanceKey string
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

func SetClientMeta(ctx context.Context, meta *ClientMeta) context.Context {
	return context.WithValue(ctx, ClientMetaContextKey{}, meta)
}

