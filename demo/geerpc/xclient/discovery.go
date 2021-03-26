package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

//负载均衡模式
type SelectMode int

const (
	//负载均衡模式-随机
	RandomSelect     SelectMode = iota
	//负载均衡模式-轮询
	RoundRobinSelect
)

//负载均衡接口
type Discovery interface {
	//刷新服务信息
	Refresh() error
	//更新服务信息
	Update(servers []string) error
	//获取服务信息，需要携带模式
	Get(mode SelectMode) (string, error)
	//获取所有服务信息
	GetAll() ([]string, error)
}

var _ Discovery = (*MultiServersDiscovery)(nil)

//负载均衡实现
//最好不要实现Get
//新建多种模式对应的服务均衡对象，持有当前对象，只实现Get
type MultiServersDiscovery struct {
	//随机对象
	r       *rand.Rand
	//服务信息修改锁
	mu      sync.RWMutex
	//所有服务信息
	servers []string
	//轮询需要用到的上一次索引
	index   int
}

func (d *MultiServersDiscovery) Refresh() error {
	return nil
}

//修改所有服务信息
func (d *MultiServersDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

//获取服务信息
func (d *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := len(d.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	//根据模式选择
	//这里写的不好，如果新增一种模式，需要修改这里的代码
	switch mode {
	case RandomSelect:
		return d.servers[d.r.Intn(n)], nil
	case RoundRobinSelect:
		s := d.servers[d.index%n]
		d.index = (d.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

//获取所有服务信息
func (d *MultiServersDiscovery) GetAll() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	//全量拷贝
	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

//新建负载均衡实现对象
func NewMultiServerDiscovery(servers []string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}
