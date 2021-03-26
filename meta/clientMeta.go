package meta

import (
	"context"
	"google.golang.org/grpc"
	"myRPC/registry/register"
)

//客户端中间件参数
const (
	//服务路由参数
	//不指定哪个节点，由均衡算法算出
	Caller_type_balance = 1
	//指定服务名称和服务id，调用指定节点
	Caller_type_one = 2
	//指定服务名称，调用所有节点
	Caller_type_all = 3
	//白名单和黑名单
)

const (
	Default_max_reconnect = 3
)

const (
	//grpc模式
	//简单模式
	Caller_mode_simple = 1
	//流模式
	Caller_mode_stream = 2
)

type ClientMetaOption func(*ClientMeta)

type ModeFunc func(interface{})

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
	//调用模式
	CallerMode int
	//流模式回调函数
	CallerModeFunc ModeFunc
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

//赋值参数函数
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

func SetCallerMode(callerMode int) ClientMetaOption {
	return func(clientMeta *ClientMeta) {
		clientMeta.CallerMode = callerMode
	}
}

func SetCallerModeFunc(callerModeFunc ModeFunc) ClientMetaOption {
	return func(clientMeta *ClientMeta) {
		clientMeta.CallerModeFunc = callerModeFunc
	}
}

func SetServiceName(serviceName string) ClientMetaOption {
	return func(clientMeta *ClientMeta) {
		clientMeta.ServiceName = serviceName
	}
}

func SetServiceMethod(serviceMethod string) ClientMetaOption {
	return func(clientMeta *ClientMeta) {
		clientMeta.ServiceMethod = serviceMethod
	}
}

type ClientMetaContextKey struct{}

//从上下文获取客户端meta
func GetClientMeta(ctx context.Context) *ClientMeta {
	meta, ok := ctx.Value(ClientMetaContextKey{}).(*ClientMeta)
	if !ok {
		meta = &ClientMeta{}
	}
	return meta
}

//保存客户端meta到上下文
func SetClientMeta(ctx context.Context, meta *ClientMeta) context.Context {
	return context.WithValue(ctx, ClientMetaContextKey{}, meta)
}

