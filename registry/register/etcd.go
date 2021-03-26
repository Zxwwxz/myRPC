package register

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"sync"
	"sync/atomic"
	"time"
)

const(
	//插件名称
	register_name_etcd = "etcd"
	//最大服务数量
	default_max_server_num = 100
	//etcd地址
	default_addr = "47.92.212.70:2379"
	//etcd路径
	default_path = "/myRpc"
	//注册超时时间
	default_timeout = 1
	//当前服务状态上报间隔
	default_report_time = 5
	//获取所有服务最新状态间隔
	default_update_time = 10
)

type EtcdRegister struct {
	//连接etcd客户端
	client *clientv3.Client
	//注册服务通道
	registerChan chan *Service
	//存储所有服务
	allServiceValue atomic.Value
	//当前服务的对象
	curService *RegisterService
	//etcd锁
	lock sync.Mutex
	//etcd地址
	addr string
	//etcd路径
	path string
	//注册超时时间
	timeout int64
	//当前服务状态上报间隔
	reportTime int
	//获取所有服务最新状态间隔
	updateTime int
}

//存储所有服务信息结构
type AllServiceInfo struct {
	allServiceMap map[string]*Service
}

//注册服务的结构
type RegisterService struct {
	//租期id
	id clientv3.LeaseID
	//续期回复通道
	reportChan <-chan *clientv3.LeaseKeepAliveResponse
	//注册服务信息
	service *Service
	//是否注册
	registered bool
}

func NewEtcdRegister(params map[interface{}]interface{})(*EtcdRegister,error) {
	addr := params["addr"].(string)
	if addr == ""{
		addr = default_addr
	}
	path := params["path"].(string)
	if path == ""{
		path = default_path
	}
	timeout := params["timeout"].(int)
	if timeout == 0{
		timeout = default_timeout
	}
	reportTime := params["report_time"].(int)
	if reportTime == 0{
		reportTime = default_report_time
	}
	updateTime := params["update_time"].(int)
	if updateTime == 0{
		updateTime = default_update_time
	}
	//创建etcd客户端
	etcdClient,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{addr},
		DialTimeout:time.Duration(timeout)*time.Second,
	})
	if err != nil {
		return nil,err
	}
	allServiceInfo := &AllServiceInfo{
		allServiceMap:make(map[string]*Service,default_max_server_num),
	}
	etcdRegister := &EtcdRegister{
		addr:addr,
		path:path,
		timeout:int64(timeout),
		reportTime:reportTime,
		updateTime:updateTime,
		client:etcdClient,
		registerChan:make(chan *Service,default_max_server_num),
	}
	etcdRegister.allServiceValue.Store(allServiceInfo)
	//开始轮询
	go etcdRegister.run()
	return etcdRegister,nil
}

//插件名字
func (e *EtcdRegister)Name()(name string){
 	return register_name_etcd
 }

//注册服务中的节点
func (e *EtcdRegister)Register(service *Service)(err error){
	//放入通道中，自然有协程会取
	select {
	case e.registerChan <- service:
	default:
		fmt.Println("Register fail")
	}
	return
}

//反注册服务中的节点
func (e *EtcdRegister)UnRegister(service *Service)(err error){
	//有租期无需反注册
	return
}

//获取服务
func (e *EtcdRegister)GetService(ctx context.Context,serviceName string)(service *Service,err error){
	//缓存中有直接返回
	localService := e.getLocalServiceInfo(serviceName)
	if localService != nil {
		return localService,nil
	}
	//没有进行etcd拉取，进行存储
	//加锁为了防止同一个服务因为没有缓存全部去请求
	e.lock.Lock()
	defer e.lock.Unlock()
	//再次检测缓存中有直接返回
	localService = e.getLocalServiceInfo(serviceName)
	if localService != nil {
		return localService,nil
	}
	path := e.getServicePath(serviceName)
	//用前缀读
	resp ,err := e.client.Get(ctx,path,clientv3.WithPrefix())
	if err != nil{
		return nil,err
	}
	remoteService := &Service{}
	remoteService.SvrNodes = []*Node{}
	for _,v := range resp.Kvs{
		tempService := &Service{}
		err := json.Unmarshal(v.Value,tempService)
		if err != nil {
			continue
		}
		remoteService.SvrName = tempService.SvrName
		remoteService.SvrType = tempService.SvrType
		for _, node := range tempService.SvrNodes {
			remoteService.SvrNodes = append(remoteService.SvrNodes,node)
		}
	}
	allServiceInfo,ok := e.allServiceValue.Load().(*AllServiceInfo)
	if ok == true && remoteService.SvrName !=  ""{
		allServiceInfo.allServiceMap[remoteService.SvrName] = remoteService
		e.allServiceValue.Store(allServiceInfo)
	}
	return remoteService,nil
}

func (e *EtcdRegister)run(){
	timer := time.NewTicker(time.Duration(e.updateTime)*time.Second)
	for {
		select {
		//有注册需要处理
		case service := <- e.registerChan:
			if e.curService == nil {
				registerService := &RegisterService{
					service:service,
					registered:false,
				}
				e.curService = registerService
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

func (e *EtcdRegister)updateAllServiceInfo()(err error)  {
	allServiceInfoOld,ok := e.allServiceValue.Load().(*AllServiceInfo)
	if ok == false || allServiceInfoOld == nil {
		return errors.New("allServiceValue Load fail")
	}
	var allServiceInfoNew = &AllServiceInfo{
		allServiceMap: make(map[string]*Service, default_max_server_num),
	}
	//旧缓存中的服务，逐一去etcd拉取更新
	for _,oldService := range allServiceInfoOld.allServiceMap{
		path := e.getServicePath(oldService.SvrName)
		resp,err := e.client.Get(context.TODO(),path,clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.allServiceMap[oldService.SvrName] = oldService
			continue
		}
		newService := &Service{
			SvrName:oldService.SvrName,
			SvrType:oldService.SvrType,
			SvrNodes: []*Node{},
		}
		for _,v := range resp.Kvs{
			tempService := &Service{}
			err := json.Unmarshal(v.Value,tempService)
			if err != nil {
				continue
			}
			for _, node := range tempService.SvrNodes {
				newService.SvrNodes = append(newService.SvrNodes,node)
			}
		}
		allServiceInfoNew.allServiceMap[newService.SvrName] = newService
	}
	e.allServiceValue.Store(allServiceInfoNew)
	return
}

func (e *EtcdRegister)checkRegisterServiceInfo() {
	if e.curService == nil {
		time.Sleep(time.Second * 1)
		return
	}
	if e.curService.registered == false {
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

func (e *EtcdRegister)registerServiceInfo()(err error) {
	//获取租期
	resp,err := e.client.Grant(context.TODO(),e.timeout)
	if err != nil {
		return
	}
	e.curService.id = resp.ID
	serviceInfo := e.curService.service
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
	e.curService.registered = true
	e.curService.reportChan = ch
	return
}

func (e *EtcdRegister)checkServiceInfo()(err error) {
	select {
	case ch := <-e.curService.reportChan:
		//说明连接断开
		if ch == nil {
			e.curService.registered = false
			return errors.New("service disconnect")
		}
	}
	return
}

func (e *EtcdRegister)getLocalServiceInfo(serviceName string)(service *Service){
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

func (e *EtcdRegister)getServicePath(serviceName string)(path string){
	return fmt.Sprintf("%s/%s",e.path,serviceName)
}

func (e *EtcdRegister)getServiceNodePath(service *Service)(path string){
	var node *Node
	for _,v := range service.SvrNodes{
		node = v
	}
	return fmt.Sprintf("%s/%s/%d",e.path,service.SvrName,node.NodeId)
}