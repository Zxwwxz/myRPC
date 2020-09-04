package meta

import (
	"context"
	"google.golang.org/grpc"
	"myRPC/registry/register"
)

const (
	//不指定哪个节点，由均衡算法算出
	Caller_type_balance = 1
	//指定服务名称和服务id，调用指定节点
	Caller_type_one = 2
	//指定服务名称，调用所有节点
	Caller_type_all = 3
)

const (
	Default_max_reconnect = 3
)

type ClientMetaOption func(*ClientMeta)

type ClientMeta struct {
	//服务端服务名
	ServiceName string
	//服务方法
	ServiceMethod string

	//客户端服务名
	ClientName string

	//调用类型
	CallerType int
	//当调用类型是3时，指定的服务id
	CallerServerId int
	//调用失败最大重连次数
	MaxReconnectNum int
	//负载均衡关键字
	BalanceKey string

	//当前节点
	CurNode *register.Node
	//历史选择节点
	RemainNodes []*register.Node
	//服务提供方的节点列表
	AllNodes []*register.Node
	//当前请求使用的连接
	Conn *grpc.ClientConn
}

func SetCallerType(callerType int) ClientMetaOption {
	return func(clientMeta *ClientMeta) {
		clientMeta.CallerType = callerType
	}
}

func SetCallerServerId(callerServerId int) ClientMetaOption {
	return func(clientMeta *ClientMeta) {
		clientMeta.CallerServerId = callerServerId
	}
}

func SetMaxReconnectNum(maxReconnectNum int) ClientMetaOption {
	return func(clientMeta *ClientMeta) {
		clientMeta.MaxReconnectNum = maxReconnectNum
	}
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

