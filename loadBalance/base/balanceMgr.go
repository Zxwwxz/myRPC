package balanceBase

import (
	"errors"
	"myRPC/loadBalance/balancer"
)

var balanceManager *BalanceManager

type BalanceManager struct {
	curBalancer balancer.BalanceInterface
}

//启动服务时初始化
func InitBalance()  {
	balanceManager = &BalanceManager{}
}

//获取全局均衡管理器
func GetBalanceMgr()*BalanceManager {
	return balanceManager
}

//获取当前负载均衡器
func GetCurBalancer()(balancer.BalanceInterface){
	if balanceManager != nil {
		return balanceManager.curBalancer
	}
	return nil
}

//每次rpc调用创建一个均衡器
func (l *BalanceManager)NewBalancer(balanceType string) (newBalancer balancer.BalanceInterface,err error) {
	switch balanceType {
	case "random":
		newBalancer = balancer.NewRandomBalance()
	case "round":
		newBalancer = balancer.NewRoundBalance()
	case "weight_random":
		newBalancer = balancer.NewWeightRandomBalance()
	}
	if newBalancer != nil {
		l.curBalancer = newBalancer
		return newBalancer,nil
	}
	return nil,errors.New("balanceType illegal")
}

func (l *BalanceManager)Stop()() {

}
