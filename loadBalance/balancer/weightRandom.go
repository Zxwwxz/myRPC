package balancer

import (
	"context"
	"errors"
	"math/rand"
	"myRPC/registry/register"
)
//加权负载均衡
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

func (r *WeightRandomBalance)SelectNode(ctx context.Context,nodes []*register.Node,params interface{})(node *register.Node,err error) {
	nodeCount := len(nodes)
	if nodeCount == 0 {
		return nil,errors.New("nodes nil")
	}
	//总权重
	weightSum := 0
	for _,v := range nodes{
		tempWeight := v.NodeWeight
		if tempWeight == 0 {tempWeight = 1}
		weightSum = weightSum + tempWeight
	}
	//落到那个区间就是哪个节点
	//节点1:100，节点2:50
	//1-150随机，那么落到100以内的概率2/3
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