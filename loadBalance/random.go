package loadBalance

import (
	"context"
	"errors"
	registryBase "myRPC/registry/base"
	"math/rand"
)

type RandomBalance struct {
	Name string
}

func NewRandomBalance() BalanceInterface {
	return &RandomBalance{
		Name:"random",
	}
}

func (r *RandomBalance)GetName()(name string) {
	return r.Name
}

func (r *RandomBalance)SelectNode(ctx context.Context,nodes map[int]*registryBase.Node)(node *registryBase.Node,err error) {
	nodeCount := len(nodes)
	if nodeCount == 0 {
		return nil,errors.New("nodes nil")
	}

	randNum := rand.Intn(nodeCount)
	return nodes[randNum],nil
}

