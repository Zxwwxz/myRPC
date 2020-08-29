package balancer

import (
	"context"
	"errors"
	"myRPC/registry/register"
)

type RoundBalance struct {
	Name string
	IndexMap map[string]int
}

func NewRoundBalance() BalanceInterface {
	return &RoundBalance{
		Name:"round",
		IndexMap:make(map[string]int),
	}
}

func (r *RoundBalance)GetName()(name string) {
	return r.Name
}

func (r *RoundBalance)SelectNode(ctx context.Context,nodes []*register.Node,params interface{})(node *register.Node,err error) {
	nodeCount := len(nodes)
	if nodeCount == 0 {
		return nil,errors.New("nodes nil")
	}
	svrName,ok := params.(string)
	if ok == false {
		return nil,errors.New("params nil")
	}
	index,ok := r.IndexMap[svrName]
	if ok == false{
		r.IndexMap[svrName] = 0
		index = 0
	}
	index = (index + 1)%nodeCount
	r.IndexMap[svrName] = index
	return nodes[index],nil
}

