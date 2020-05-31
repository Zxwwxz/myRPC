package balanceBase

import (
	"context"
	registryBase "github.com/ibinarytree/koala/registry/base"
)

type BalanceInterface interface {
	GetName()(name string)
	SelectNode(ctx context.Context,nodes []*registryBase.Node)(node *registryBase.Node,err error)
}
