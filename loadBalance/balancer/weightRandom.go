package balancer

import (
	"context"
	"errors"
	"math/rand"
)

type WeightRandomBalance struct {
	Name string
}

func NewWeightRandomBalance() BalanceInterface {
	return &WeightRandomBalance{
		Name:"weight_random",
	}
}

func (r *WeightRandomBalance)GetName()(name string) {
	return r.Name
}

func (r *WeightRandomBalance)SelectNode(ctx context.Context,nodes []*Node,params interface{})(node *Node,err error) {
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
	randNum := rand.Intn(weightSum)
	for i,v := range nodes{
		tempWeight := v.NodeWeight
		if tempWeight == 0 {tempWeight = 1}
		randNum = randNum - tempWeight
		if randNum < 0 {
			return nodes[i],nil
		}
	}
	return nil,errors.New("not found node")
}