package balancer

import (
	"context"
)

type BalanceInterface interface {
	GetName()(name string)
	SelectNode(ctx context.Context,nodes []*Node,params interface{})(node *Node,err error)
}
