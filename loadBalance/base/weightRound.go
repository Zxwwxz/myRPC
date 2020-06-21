package balanceBase

import (
	"context"
	"errors"
	registryBase "myRPC/registry/base"
)

type WeightRoundBalance struct {
	Name string
	Index int
}

func NewWeightRoundBalance() BalanceInterface {
	return &WeightRoundBalance{
		Name:"weight_round",
		Index:0,
	}
}

func (r *WeightRoundBalance)GetName()(name string) {
	return r.Name
}

func (r *WeightRoundBalance)SelectNode(ctx context.Context,nodes map[int]*registryBase.Node)(node *registryBase.Node,err error) {
	nodeCount := len(nodes)
	if nodeCount == 0 {
		return nil,errors.New("nodes nil")
	}
	weightSum := 0
	for _,v := range nodes{
		tempWeight := v.NodeWeight
		if tempWeight == 0 {tempWeight = 1}
		weightSum = weightSum + tempWeight
	}
	index := (r.Index + 1)%weightSum
	for i,v := range nodes{
		tempWeight := v.NodeWeight
		if tempWeight == 0 {tempWeight = 1}
		index = index - tempWeight
		if index < 0 {
			r.Index = index
			return nodes[i],nil
		}
	}
	return nodes[r.Index],nil
}

