package xclient

import (
	"log"
	"net/http"
	"strings"
	"time"
)

//带有获取注册中心型负载均衡
type GeeRegistryDiscovery struct {
	//里面实现了服务均衡
	*MultiServersDiscovery
	//注册中心地址
	registry   string
	//服务列表过期时间，超过这个时间需要前去拉取
	timeout    time.Duration
	//上一次去注册中心拉取的时间
	lastUpdate time.Time
}

//默认去注册中心拉取的时间间隔
const defaultUpdateTimeout = time.Second * 10

//修改服务列表
func (d *GeeRegistryDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

//前往注册中心拉取服务列表
func (d *GeeRegistryDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	//没超时，不用更新
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}
	//发送http到注册中心拉取，更新本地
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Println("rpc registry refresh err:", err)
		return err
	}
	servers := strings.Split(resp.Header.Get("X-Geerpc-Servers"), ",")
	d.servers = make([]string, 0, len(servers))
	for _, server := range servers {
		if strings.TrimSpace(server) != "" {
			d.servers = append(d.servers, strings.TrimSpace(server))
		}
	}
	d.lastUpdate = time.Now()
	return nil
}

//先进行刷新，然后直接调用负载均衡对象的方法即可
func (d *GeeRegistryDiscovery) Get(mode SelectMode) (string, error) {
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}

func (d *GeeRegistryDiscovery) GetAll() ([]string, error) {
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}

//创建对象
func NewGeeRegistryDiscovery(registerAddr string, timeout time.Duration) *GeeRegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}
	d := &GeeRegistryDiscovery{
		MultiServersDiscovery: NewMultiServerDiscovery(make([]string, 0)),
		registry:              registerAddr,
		timeout:               timeout,
	}
	return d
}
