package registryBase

import (
	"errors"
	"myRPC/registry/register"
)

var (
	registryManager *RegistryManager
)

//插件管理器
type RegistryManager struct {
	//当前注册器对象
	curRegister register.RegisterInterface
}

func InitRegistry() {
	registryManager = &RegistryManager{}
}

func GetRegistryManager()(*RegistryManager){
	return registryManager
}

func GetRegister()(register.RegisterInterface){
	if registryManager != nil {
		return registryManager.curRegister
	}
	return nil
}

//新建插件
func (r *RegistryManager)NewRegister(registryType string,params map[interface{}]interface{})(newRegister register.RegisterInterface,err error){
	switch registryType {
	case "etcd":
		newRegister,err = register.NewEtcdRegister(params)
	}
	if newRegister != nil {
		r.curRegister = newRegister
		return newRegister,nil
	}
	return nil,errors.New("registryType illegal")
}

func (r *RegistryManager)Stop()(){

}

//注册当前服务
func (r *RegistryManager)RegisterServer(server *register.Service)(err error){
	if r.curRegister != nil {
		return r.curRegister.Register(server)
	}
	return errors.New("curRegister nil")
}
