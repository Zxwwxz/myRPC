package base

import "time"

//当前服务注册相关参数
type RegisterOptions struct {
	//注册地址
	RegisterAddrs []string
	//注册超时
	RegisterTimeOut time.Duration
	//注册路径
	RegisterPath string
	//心跳超时，租期续约
	HeartTimeOut int64
}

//当前服务注册函数类型
type RegisterOptionFunc func(registerOptions *RegisterOptions)

//具体设置值函数
func SetRegisterAddrs(registerAddrs []string) RegisterOptionFunc{
	return func(registerOptions *RegisterOptions) {
		registerOptions.RegisterAddrs = registerAddrs
	}
}

func SetRegisterTimeOut(registerTimeOut time.Duration) RegisterOptionFunc{
	return func(registerOptions *RegisterOptions) {
		registerOptions.RegisterTimeOut = registerTimeOut
	}
}

func SetRegisterPath(registerPath string) RegisterOptionFunc{
	return func(registerOptions *RegisterOptions) {
		registerOptions.RegisterPath = registerPath
	}
}

func SetHeartTimeOut(heartTimeOut int64) RegisterOptionFunc{
	return func(registerOptions *RegisterOptions) {
		registerOptions.HeartTimeOut = heartTimeOut
	}
}

