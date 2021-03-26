package rpcConst

import "fmt"

type ClientError struct {
	ErrCode    int
	Message string
}

func (k *ClientError) Error() string {
	return fmt.Sprintf("client error, code:%d message:%v", k.ErrCode, k.Message)
}

func (k *ClientError) Code() int {
	return k.ErrCode
}

var (
	//客户端初始化失败
	ClientInitFailed = &ClientError{
		ErrCode:    1,
		Message: "client init failed",
	}
	//连接失败
	ConnFailed = &ClientError{
		ErrCode:    2,
		Message: "connect failed",
	}
	//未找到节点
	NotFoundNode = &ClientError{
		ErrCode:    3,
		Message: "not found node",
	}
	//所有节点访问失败
	AllNodeFailed = &ClientError{
		ErrCode:    4,
		Message: "all node failed",
	}
	//最大重连失败
	MaxReconnectFailed = &ClientError{
		ErrCode:    5,
		Message: "max reconnect failed",
	}
	//客户端限流
	ClientLimit = &ClientError{
		ErrCode:    6,
		Message: "client rate limited",
	}
	//返回值类型错误
	ClientReturnIllegal = &ClientError{
		ErrCode:    7,
		Message: "client return type illegal",
	}
)
