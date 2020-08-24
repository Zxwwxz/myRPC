package balanceBase

import (
	"errors"
	"myRPC/loadBalance/balancer"
)

var balanceManager *BalanceManager

type BalanceManager struct {}

//启动服务时初始化
func InitBalance()  {
	balanceManager = &BalanceManager{}
}

//获取全局均衡管理器
func GetBalanceMgr()*BalanceManager {
	return balanceManager
}

//每次rpc调用创建一个均衡器
func (l *BalanceManager)NewBalancer(balanceType string) (balancer.BalanceInterface,error) {
	switch balanceType {
	case "random":
		return balancer.NewRandomBalance(),nil
	case "round":
		return balancer.NewRoundBalance(),nil
	case "weight_random":
		return balancer.NewWeightRandomBalance(),nil
	}
	return nil,errors.New("balanceType illegal")
}
