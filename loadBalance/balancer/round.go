package balancer

import (
	"context"
	"errors"
)

type RoundBalance struct {
	Name string
	IndexMap map[int]int
}

func NewRoundBalance() BalanceInterface {
	return &RoundBalance{
		Name:"round",
		IndexMap:make(map[int]int),
	}
}

func (r *RoundBalance)GetName()(name string) {
	return r.Name
}

func (r *RoundBalance)SelectNode(ctx context.Context,nodes []*Node,params interface{})(node *Node,err error) {
	nodeCount := len(nodes)
	if nodeCount == 0 {
		return nil,errors.New("nodes nil")
	}
	svrType,ok := params.(int)
	if ok == false {
		return nil,errors.New("params nil")
	}
	index,ok := r.IndexMap[svrType]
	if ok == false{
		r.IndexMap[svrType] = 0
		index = 0
	}
	index = (index + 1)%nodeCount
	r.IndexMap[svrType] = index
	return nodes[index],nil
}

