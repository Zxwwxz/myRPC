package rpcConst

import "fmt"

type ClientError struct {
	Code    int
	Message string
}

func (k *ClientError) Error() string {
	return fmt.Sprintf("client error, code:%d message:%v", k.Code, k.Message)
}

var (
	//连接失败
	ConnFailed = &ClientError{
		Code:    1,
		Message: "connect failed",
	}
	//未找到节点
	NotFoundNode = &ClientError{
		Code:    2,
		Message: "not found node",
	}
	//所有节点访问失败
	AllNodeFailed = &ClientError{
		Code:    3,
		Message: "all node failed",
	}
)
