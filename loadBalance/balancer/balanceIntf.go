package balancer

import (
	"context"
	"myRPC/registry/register"
)

type BalanceInterface interface {
	GetName()(name string)
	SelectNode(ctx context.Context,nodes []*register.Node,params interface{})(node *register.Node,err error)
}
