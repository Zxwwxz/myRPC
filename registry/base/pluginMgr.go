package base

import (
	"context"
	"errors"
	"sync"
)

//插件管理器
type pluginManager struct {
	plugins map[string]RegistryPlugin
	lock sync.Mutex
}

//插件管理器全局对象
var PluginManager = &pluginManager{
	plugins:make(map[string]RegistryPlugin),
}

//注册插件
func (p *pluginManager)RegisterPlugin(registryPlugin RegistryPlugin)(err error){
	p.lock.Lock()
	defer p.lock.Unlock()
	if _,ok := p.plugins[registryPlugin.Name()];ok == true {
		return errors.New("already register plugin")
	}
	p.plugins[registryPlugin.Name()] = registryPlugin
	return
}

//初始化插件
func (p *pluginManager)InitPlugin(ctx context.Context, name string, registerOptionFuncs ...RegisterOptionFunc)(plugin RegistryPlugin, err error){
	p.lock.Lock()
	defer p.lock.Unlock()
	if plugin,ok := p.plugins[name];ok == true {
		err := plugin.Init(ctx,registerOptionFuncs...)
		return plugin,err
	}else{
		return plugin,errors.New("not found plugin")
	}
}
