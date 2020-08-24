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

//注册插件
func (r *RegistryManager)NewRegister(registryType string,params map[interface{}]interface{})(newRegister register.RegisterInterface,err error){
	switch registryType {
	case "etcd":
		newRegister,err = register.NewEtcdRegister(params)
	}
	if newRegister != nil {
		r.curRegister = newRegister
	}
	return nil,errors.New("limiterType illegal")
}

func (r *RegistryManager)RegisterServer(server *register.Service)(err error){
	if r.curRegister != nil {
		return r.curRegister.Register(server)
	}
	return errors.New("curRegister nil")
}
