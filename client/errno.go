package client

import "fmt"

type ClientError struct {
	Code    int
	Message string
}

func (k *ClientError) Error() string {
	return fmt.Sprintf("client error, code:%d message:%v", k.Code, k.Message)
}

var (
	NotHaveInstance = &ClientError{
		Code:    1,
		Message: "not have instance",
	}
	ConnFailed = &ClientError{
		Code:    2,
		Message: "connect failed",
	}
	InvalidNode = &ClientError{
		Code:    3,
		Message: "invalid node",
	}
	AllNodeFailed = &ClientError{
		Code:    4,
		Message: "all node failed",
	}
)
