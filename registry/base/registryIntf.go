package base

import "context"

//注册插件接口
type RegistryPlugin interface {
	//插件名字
	Name()(name string)
	//初始化
	Init(ctx context.Context,registerOptionFuncs ...RegisterOptionFunc)(err error)
	//注册服务中的节点
	Register(ctx context.Context,service *Service)(err error)
	//反注册服务中的节点
	UnRegister(ctx context.Context,service *Service)(err error)
	//获取服务
	GetService(ctx context.Context,serviceName string)(service *Service,err error)
}
