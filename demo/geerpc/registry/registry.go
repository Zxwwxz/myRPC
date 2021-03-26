package registry

import (
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

//注册中心，管理所有可用的服务信息
type GeeRegistry struct {
	//单个服务的可用有效时间，需要服务每隔一段时间进行心跳保活，超过了时间没收到心跳该服务会失效
	timeout time.Duration
	//注册中心的服务修改锁
	mu      sync.Mutex
	//所有可用服务信息，key是ip:port，value是服务信息
	servers map[string]*ServerItem
}

//单个服务的信息，可以添加更多的服务相关信息
type ServerItem struct {
	//服务ip:port
	Addr  string
	//服务最近一次心跳包保活开始时间
	start time.Time
}

const (
	//服务注册的http路径
	defaultPath    = "/_geerpc_/registry"
	//默认服务可用有效时间
	defaultTimeout = time.Minute * 5
)

//新建注册中心
func New(timeout time.Duration) *GeeRegistry {
	return &GeeRegistry{
		servers: make(map[string]*ServerItem),
		timeout: timeout,
	}
}

//默认全局注册中心
var DefaultGeeRegister = New(defaultTimeout)

//修改服务信息到注册中心
func (r *GeeRegistry) putServer(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	//获取服务信息
	s := r.servers[addr]
	if s == nil {
		//新增信息
		r.servers[addr] = &ServerItem{Addr: addr, start: time.Now()}
	} else {
		//修改保活开始时间
		s.start = time.Now() // if exists, update start time to keep alive
	}
}

//检测所有服务是否存活，并返回可用的服务
func (r *GeeRegistry) aliveServers() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var alive []string
	for addr, s := range r.servers {
		//如果是永久存活，或者还在时间范围内
		if r.timeout == 0 || s.start.Add(r.timeout).After(time.Now()) {
			alive = append(alive, addr)
		} else {
			//已经失效，进行删除
			delete(r.servers, addr)
		}
	}
	sort.Strings(alive)
	return alive
}

//修改或者拉取服务配置
func (r *GeeRegistry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		//添加所有存活服务到http头中返回
		w.Header().Set("X-Geerpc-Servers", strings.Join(r.aliveServers(), ","))
	case "POST":
		//将http头部服务信息取出，修改注册中心服务信息
		addr := req.Header.Get("X-Geerpc-Server")
		if addr == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.putServer(addr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//开启http服务
func (r *GeeRegistry) HandleHTTP(registryPath string) {
	http.Handle(registryPath, r)
}

func HandleHTTP() {
	DefaultGeeRegister.HandleHTTP(defaultPath)
}

//服务定时发送心跳包，前面是注册中心http地址，后面是当前服务地址信息
func Heartbeat(registry, addr string, duration time.Duration) {
	if duration == 0 {
		duration = defaultTimeout - time.Duration(1)*time.Minute
	}
	var err error
	err = sendHeartbeat(registry, addr)
	//开启定时器，发送心跳包，发送时间间隔要短于服务过期时间
	go func() {
		t := time.NewTicker(duration)
		for err == nil {
			<-t.C
			err = sendHeartbeat(registry, addr)
		}
	}()
}

func sendHeartbeat(registry, addr string) error {
	httpClient := &http.Client{}
	//发送post请求，信息添加到头部
	req, _ := http.NewRequest("POST", registry, nil)
	req.Header.Set("X-Geerpc-Server", addr)
	if _, err := httpClient.Do(req); err != nil {
		return err
	}
	return nil
}
