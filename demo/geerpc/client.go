// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geerpc

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

//每次rpc调用结构体
type Call struct {
	//当前客户端调用顺序
	Seq           uint64
	//调用服务端接口名
	ServiceMethod string
	//请求参数
	Args          interface{}
	//返回参数
	Reply         interface{}
	//调用错误
	Error         error
	//调用完成通知
	Done          chan *Call
}

//进行通知，将当前调用结构体传给等待的通道
func (call *Call) done() {
	call.Done <- call
}

//客户端对象
type Client struct {
	//编解码对象
	cc       codec.Codec
	//客户端参数
	opt      *Option
	//发送同步锁
	sending  sync.Mutex
	//复用的每次rpc调用头
	header   codec.Header
	//当前对象锁
	mu       sync.Mutex
	//客户端请求序列
	seq      uint64
	//还没收到确认的rpc请求
	pending  map[uint64]*Call
	//客户端是否已经关闭，一般是因为主动关闭
	closing  bool
	//客户端是否已经休眠，一般是因为有错误
	shutdown bool
}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connection is shut down")

//关闭客户端
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

//是否可用
func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

//注册调用
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	//将调用保存到未确认中
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

//移除调用
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

//通知所有未确认调用连接已经断开
func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

//发送rpc
func (client *Client) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()
	//先注册
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}
	//封装请求头
	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""
	//通过编解码对象发送到链路中
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		//失败了
		call := client.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

//接受服务端回复
func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		//解析头
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}
		//移除，表示已经收到调用
		call := client.removeCall(h.Seq)
		switch {
		//可能是之前有异常，每保存对象
		case call == nil:
			err = client.cc.ReadBody(nil)
			//调用有错误
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			call.done()
			//正常返回
		default:
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	//退出循环表示有异常
	client.terminateCalls(err)
}

//异步调用，调用后即可得到call对象
//要获取结果自己开启协程，等待call.Done返回值
//不用获取结果直接不管
func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	client.send(call)
	return call
}

//同步调用，直到服务端返回结果才调用结束
func (client *Client) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	call := client.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	//上下文可以设置超时时间，超时了不等待直接返回
	case <-ctx.Done():
		client.removeCall(call.Seq)
		return errors.New("rpc client: call failed: " + ctx.Err().Error())
	case call := <-call.Done:
		return call.Error
	}
}

//解析参数
func parseOptions(opts ...*Option) (*Option, error) {
	if len(opts) == 0 || opts[0] == nil {
		return DefaultOption, nil
	}
	if len(opts) != 1 {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = DefaultOption.CodecType
	}
	return opt, nil
}

//创建Tcp客户端对象
func NewClient(conn net.Conn, opt *Option) (*Client, error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error:", err)
		return nil, err
	}
	//连接建立后先发送参数
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error: ", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

//创建携带编解码对象的客户端对象
func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		seq:     1, // seq starts with 1, 0 means invalid call
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
	}
	//开启协程循环接收服务端消息
	go client.receive()
	return client
}

//客户端创建结果
type clientResult struct {
	client *Client
	err    error
}

//建立连接函数
type newClientFunc func(conn net.Conn, opt *Option) (client *Client, err error)

//建立统一连接客户端，包装了连接超时代码
func dialTimeout(f newClientFunc, network, address string, opts ...*Option) (client *Client, err error) {
	//解析参数
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	//建立连接对象，设置连接超时时间
	conn, err := net.DialTimeout(network, address, opt.ConnectTimeout)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()
	ch := make(chan clientResult)
	go func() {
		//调用对应的tcp或http构建函数得到客户端
		client, err := f(conn, opt)
		ch <- clientResult{client: client, err: err}
	}()
	//没有设置超时，就一直等待
	if opt.ConnectTimeout == 0 {
		result := <-ch
		return result.client, result.err
	}
	select {
	//有设置超时
	case <-time.After(opt.ConnectTimeout):
		return nil, fmt.Errorf("rpc client: connect timeout: expect within %s", opt.ConnectTimeout)
	case result := <-ch:
		return result.client, result.err
	}
}

//建立连接tcp连接
func Dial(network, address string, opts ...*Option) (*Client, error) {
	return dialTimeout(NewClient, network, address, opts...)
}

//建立http连接客户端
func NewHTTPClient(conn net.Conn, opt *Option) (*Client, error) {
	//发送http的Connect方法
	_, _ = io.WriteString(conn, fmt.Sprintf("CONNECT %s HTTP/1.0\n\n", defaultRPCPath))
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	//得到服务端200的结果，将http连接转成tcp连接
	if err == nil && resp.Status == connected {
		return NewClient(conn, opt)
	}
	if err == nil {
		err = errors.New("unexpected HTTP response: " + resp.Status)
	}
	return nil, err
}

//建立http连接
func DialHTTP(network, address string, opts ...*Option) (*Client, error) {
	return dialTimeout(NewHTTPClient, network, address, opts...)
}

//建立连接统一入口
//http@10.0.0.1:7001, tcp@10.0.0.1:9999, unix@/tmp/geerpc.sock
func XDial(rpcAddr string, opts ...*Option) (*Client, error) {
	parts := strings.Split(rpcAddr, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("rpc client err: wrong format '%s', expect protocol@addr", rpcAddr)
	}
	//协议和地址
	protocol, addr := parts[0], parts[1]
	switch protocol {
	//建立http连接
	case "http":
		return DialHTTP("tcp", addr, opts...)
		//建立tcp连接
	default:
		return Dial(protocol, addr, opts...)
	}
}
