package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"myRPC/registry/base"
	"sync"
	"sync/atomic"
	"time"
)

const(
	//插件名称
	EtcdPluginName = "etcd"
	//更新服务时间间隔
	UpdateAllServiceTime = time.Second * 5
	//最大服务数量
	MaxServiceNum = 100
	//最大注册服务数量
	MaxRegisterServiceNum = 100
)

type EtcdPlugin struct {
	//插件参数
	options *registryBase.RegisterOptions
	//连接etcd客户端
	client *clientv3.Client
	//注册服务通道
	registerChan chan *registryBase.Service
	//存储所有服务
	allServiceValue atomic.Value
	//注册服务的对象
	registerService *RegisterService
	//etcd锁
	lock sync.Mutex
}

//存储所有服务信息结构
type AllServiceInfo struct {
	allServiceMap map[string]*registryBase.Service
}

//注册服务的结构
type RegisterService struct {
	//租期id
	id clientv3.LeaseID
	//续期回复通道
	heartChan <-chan *clientv3.LeaseKeepAliveResponse
	//注册服务信息
	service *registryBase.Service
	//是否注册
	registered bool
}

var EtcdPluginObj = &EtcdPlugin{}

func init()  {
	//进行注册到插件管理器
	err := registryBase.PluginManager.RegisterPlugin(EtcdPluginObj)
	if err != nil {
		fmt.Println("RegisterPlugin err:",err)
	}
	//开启轮询
	go EtcdPluginObj.run()
}

//插件名字
func (e *EtcdPlugin)Name()(name string){
 	return EtcdPluginName
 }
//初始化
func (e *EtcdPlugin)Init(ctx context.Context,registerOptionFuncs ...registryBase.RegisterOptionFunc)(err error){
	//获得参数
	e.options = &registryBase.RegisterOptions{}
	for _,v := range registerOptionFuncs{
		v(e.options)
	}

	//初始化连接对象
	e.client,err = clientv3.New(clientv3.Config{
		Endpoints:e.options.RegisterAddrs,
		DialTimeout:e.options.RegisterTimeOut,
	})

	if err != nil {
		fmt.Println("clientv3 New err:",err)
	}

	//存储初始所有服务信息
	allServiceInfo := &AllServiceInfo{
		allServiceMap:make(map[string]*registryBase.Service,MaxServiceNum),
	}
	e.allServiceValue.Store(allServiceInfo)

	e.registerChan = make(chan *registryBase.Service,MaxRegisterServiceNum)
	return
}
//注册服务中的节点
func (e *EtcdPlugin)Register(ctx context.Context,service *registryBase.Service)(err error){
	//放入通道中，自然有协程会取
	select {
	case e.registerChan <- service:
	default:
		fmt.Println("Register fail")
	}
	return
}
//反注册服务中的节点
func (e *EtcdPlugin)UnRegister(ctx context.Context,service *registryBase.Service)(err error){
	//有租期无需反注册
	return
}
//获取服务
func (e *EtcdPlugin)GetService(ctx context.Context,serviceName string)(service *registryBase.Service,err error){
	//缓存中有直接返回
	service = e.getLocalServiceInfo(serviceName)
	if service != nil {
		return
	}
	//没有进行etcd拉取，进行存储
	//加锁为了防止同一个服务因为没有缓存全部去请求
	e.lock.Lock()
	defer e.lock.Unlock()
	//再次检测缓存中有直接返回
	service = e.getLocalServiceInfo(serviceName)
	if service != nil {
		return
	}
	path := e.getServicePath(serviceName)
	//用前缀读
	resp ,err := e.client.Get(ctx,path,clientv3.WithPrefix())
	if err != nil{
		fmt.Println("etcd get err:",err)
	}
	service = &registryBase.Service{}
	for _,v := range resp.Kvs{
		tempService := &registryBase.Service{}
		err := json.Unmarshal(v.Value,tempService)
		if err != nil {
			continue
		}
		service.SvrName = tempService.SvrName
		service.SvrType = tempService.SvrType
		service.SvrNodes = make(map[int]*registryBase.Node)
		for _, node := range tempService.SvrNodes {
			service.SvrNodes[node.NodeId] = node
		}
	}
	allServiceInfo,_ := e.allServiceValue.Load().(*AllServiceInfo)
	if service.SvrName !=  ""{
		allServiceInfo.allServiceMap[service.SvrName] = service
		e.allServiceValue.Store(allServiceInfo)
	}
	return service,nil
}

func (e *EtcdPlugin)run(){
	timer := time.NewTicker(UpdateAllServiceTime)
	for {
		select {
		//有注册需要处理
		case service := <- e.registerChan:
			if e.registerService == nil {
				registerService := &RegisterService{
					service:service,
					registered:false,
				}
				e.registerService = registerService
			}
		//需要更新一下所有服务最新信息
		case <-timer.C:
			err := e.updateAllServiceInfo()
			if err != nil {
				fmt.Println("updateAllServiceInfo err:",err)
			}
		//检查当前服务是需要注册还是续租
		default:
			e.checkRegisterServiceInfo()
		}
	}
}

func (e *EtcdPlugin)updateAllServiceInfo()(err error)  {
	allServiceInfoOld,ok := e.allServiceValue.Load().(*AllServiceInfo)
	if ok == false || allServiceInfoOld == nil {
		return errors.New("allServiceValue Load fail")
	}
	var allServiceInfoNew = &AllServiceInfo{
		allServiceMap: make(map[string]*registryBase.Service, MaxServiceNum),
	}
	//旧缓存中的服务，逐一去etcd拉取更新
	for _,oldService := range allServiceInfoOld.allServiceMap{
		path := e.getServicePath(oldService.SvrName)
		resp,err := e.client.Get(context.TODO(),path,clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.allServiceMap[oldService.SvrName]=oldService
			continue
		}
		newService := &registryBase.Service{
			SvrName:oldService.SvrName,
			SvrType:oldService.SvrType,
			SvrNodes: map[int]*registryBase.Node{},
		}
		for _,v := range resp.Kvs{
			tempService := &registryBase.Service{}
			err := json.Unmarshal(v.Value,tempService)
			if err != nil {
				continue
			}
			for _, node := range tempService.SvrNodes {
				newService.SvrNodes[node.NodeId] = node
			}
		}
		allServiceInfoNew.allServiceMap[newService.SvrName]=newService
	}
	e.allServiceValue.Store(allServiceInfoNew)
	return
}

func (e *EtcdPlugin)checkRegisterServiceInfo() {
	if e.registerService == nil {
		time.Sleep(time.Second * 1)
		return
	}
	if e.registerService.registered == false {
		//注册
		err := e.registerServiceInfo()
		if err != nil {
			fmt.Println("registerServiceInfo err:",err)
		}
		time.Sleep(time.Second * 1)
		return
	}else{
		//检测
		err := e.checkServiceInfo()
		if err != nil {
			fmt.Println("checkServiceInfo err:",err)
		}
		return
	}
}

func (e *EtcdPlugin)registerServiceInfo()(err error) {
	//获取租期
	resp,err := e.client.Grant(context.TODO(),e.options.HeartTimeOut)
	if err != nil {
		return
	}
	e.registerService.id = resp.ID
	serviceInfo := e.registerService.service
	serviceInfoJson,err := json.Marshal(serviceInfo)
	if err != nil {
		return
	}
	key := e.getServiceNodePath(serviceInfo)
	//存储当前节点信息到etcd
	_,err = e.client.Put(context.TODO(),key,string(serviceInfoJson),clientv3.WithLease(resp.ID))
	if err != nil {
		return
	}
	//自动续租
	ch,err := e.client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		return
	}
	e.registerService.registered = true
	e.registerService.heartChan = ch
	return
}

func (e *EtcdPlugin)checkServiceInfo()(err error) {
	select {
	case ch := <-e.registerService.heartChan:
		//说明连接断开
		if ch == nil {
			e.registerService.registered = false
			return errors.New("service disconnect")
		}
	}
	return
}

func (e *EtcdPlugin)getLocalServiceInfo(serviceName string)(service *registryBase.Service){
	allServiceInfo,ok := e.allServiceValue.Load().(*AllServiceInfo)
	if ok == false || allServiceInfo == nil {
		return nil
	}
	service,ok = allServiceInfo.allServiceMap[serviceName]
	if ok == false {
		return nil
	}
	return
}

func (e *EtcdPlugin)getServicePath(serviceName string)(path string){
	return fmt.Sprintf("%s/%s",e.options.RegisterPath,serviceName)
}

func (e *EtcdPlugin)getServiceNodePath(service *registryBase.Service)(path string){
	var node *registryBase.Node
	for _,v := range service.SvrNodes{
		node = v
	}
	return fmt.Sprintf("%s/%s/%d",e.options.RegisterPath,service.SvrName,node.NodeId)
}