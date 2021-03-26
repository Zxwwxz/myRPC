package balancer

import (
	"context"
	"errors"
	"myRPC/registry/register"
	"sync"
)

type RoundBalance struct {
	Name string
	//不同服务对应的index
	IndexMap map[string]int
	//因为不同服务通用一个负载均衡器，所以要加锁查看修改index
	lock sync.Mutex
}

func NewRoundBalance() BalanceInterface {
	return &RoundBalance{
		Name:"round",
		IndexMap:make(map[string]int),
		lock:sync.Mutex{},
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
	r.lock.Lock()
	defer r.lock.Unlock()
	index,ok := r.IndexMap[svrName]
	if ok == false{
		r.IndexMap[svrName] = 0
		index = 0
	}
	//+1取余
	index = (index + 1)%nodeCount
	r.IndexMap[svrName] = index
	return nodes[index],nil
}

