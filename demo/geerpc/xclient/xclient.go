package xclient

import (
	"context"
	. "geerpc"
	"io"
	"reflect"
	"sync"
)

//客户端对象封装类
type XClient struct {
	//负载均衡
	d       Discovery
	//服务均衡模式
	mode    SelectMode
	//客户端参数
	opt     *Option
	//客户端对象锁
	mu      sync.Mutex
	//客户端缓存对象
	clients map[string]*Client
}

var _ io.Closer = (*XClient)(nil)

//新建客户端对象封装类
func NewXClient(d Discovery, mode SelectMode, opt *Option) *XClient {
	return &XClient{d: d, mode: mode, opt: opt, clients: make(map[string]*Client)}
}

//关闭客户端对象封装类
func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	//将所有客户端对象关闭
	for key, client := range xc.clients {
		_ = client.Close()
		delete(xc.clients, key)
	}
	return nil
}

//客户端对象获取
func (xc *XClient) dial(rpcAddr string) (*Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	//先看缓存中有没有，而且要可用的，不可用就清理掉
	client, ok := xc.clients[rpcAddr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(xc.clients, rpcAddr)
		client = nil
	}
	//没有就建立一下连接
	if client == nil {
		var err error
		client, err = XDial(rpcAddr, xc.opt)
		if err != nil {
			return nil, err
		}
		//建立了保存，下次使用
		xc.clients[rpcAddr] = client
	}
	return client, nil
}

//客户端已知服务器地址调用rpc
func (xc *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	//先获取对象
	client, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}
	//再进行发送
	return client.Call(ctx, serviceMethod, args, reply)
}

//客户端调用rpc
func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	//先通过负载均衡获取单一服务器地址
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	//调用rpc
	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)
}

//客户端广播rpc
func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	//先通过负载均衡获取所有服务器地址
	servers, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	//控制全部服务发送完成
	var wg sync.WaitGroup
	//对多个返回值加锁
	var mu sync.Mutex
	//接受其中一个返回值的错误
	var e error
	//接受其中一个返回值的结果
	replyDone := reply == nil
	//开启可取消上下文
	ctx, cancel := context.WithCancel(ctx)
	//循环遍历发送到每个服务中
	for _, rpcAddr := range servers {
		wg.Add(1)
		go func(rpcAddr string) {
			defer wg.Done()
			var clonedReply interface{}
			if reply != nil {
				clonedReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}
			//rpc调用
			err := xc.call(rpcAddr, ctx, serviceMethod, args, clonedReply)
			mu.Lock()
			if err != nil && e == nil {
				e = err
				//有一个错误了，取消等待结果
				cancel()
			}
			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(clonedReply).Elem())
				replyDone = true
			}
			mu.Unlock()
		}(rpcAddr)
	}
	wg.Wait()
	return e
}
