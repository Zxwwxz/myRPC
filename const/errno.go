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
	//没有对象
	NotHaveInstance = &ClientError{
		Code:    1,
		Message: "not have instance",
	}
	//连接失败
	ConnFailed = &ClientError{
		Code:    2,
		Message: "connect failed",
	}
	//节点无法访问
	InvalidNode = &ClientError{
		Code:    3,
		Message: "invalid node",
	}
	//所有节点访问失败
	AllNodeFailed = &ClientError{
		Code:    4,
		Message: "all node failed",
	}
)
