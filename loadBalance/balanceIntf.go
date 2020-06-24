package loadBalance

import (
	"context"
	registryBase "myRPC/registry/base"
)

type BalanceInterface interface {
	GetName()(name string)
	SelectNode(ctx context.Context,nodes map[int]*registryBase.Node)(node *registryBase.Node,err error)
}
