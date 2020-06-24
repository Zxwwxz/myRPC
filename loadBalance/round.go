package loadBalance

import (
	"context"
	"errors"
	registryBase "myRPC/registry/base"
)

type RoundBalance struct {
	Name string
	Index int
}

func NewRoundBalance() BalanceInterface {
	return &RoundBalance{
		Name:"round",
		Index:0,
	}
}

func (r *RoundBalance)GetName()(name string) {
	return r.Name
}

func (r *RoundBalance)SelectNode(ctx context.Context,nodes map[int]*registryBase.Node)(node *registryBase.Node,err error) {
	nodeCount := len(nodes)
	if nodeCount == 0 {
		return nil,errors.New("nodes nil")
	}
	r.Index = (r.Index + 1)%nodeCount
	return nodes[r.Index],nil
}

