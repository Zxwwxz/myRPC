package register

import "context"

//注册插件接口
type RegisterInterface interface {
	//插件名字
	Name()(name string)
	//注册服务中的节点
	Register(service *Service)(err error)
	//反注册服务中的节点
	UnRegister(service *Service)(err error)
	//获取服务
	GetService(ctx context.Context,serviceName string)(service *Service,err error)
}
