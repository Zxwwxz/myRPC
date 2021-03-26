package balancer

import (
	"context"
	"errors"
	"math/rand"
	"myRPC/registry/register"
)

//随机选择
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

func (r *RandomBalance)SelectNode(ctx context.Context,nodes []*register.Node,params interface{})(node *register.Node,err error) {
	nodeCount := len(nodes)
	if nodeCount == 0 {
		return nil,errors.New("nodes nil")
	}
	//随机即可
	randNum := rand.Intn(nodeCount)
	return nodes[randNum],nil
}

