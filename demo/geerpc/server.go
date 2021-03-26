// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

//校验值
const MagicNumber = 0x3bef5c

//客户端建立连接后传递的参数
type Option struct {
	//校验值
	MagicNumber    int
	//编解码类型
	CodecType      codec.Type
	//连接超时
	ConnectTimeout time.Duration
	//处理超时
	HandleTimeout  time.Duration
}

//默认连接
var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: time.Second * 10,
}

//服务端处理对象
type Server struct {
	serviceMap sync.Map
}

//新建对象
func NewServer() *Server {
	return &Server{}
}

//默认对象
var DefaultServer = NewServer()

//单个连接的处理
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option
	//解析客户端参数
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	//校验值
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	//根据解析类型，得到对应的解析方法
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn), &opt)
}

var invalidRequest = struct{}{}

func (server *Server) serveCodec(cc codec.Codec, opt *Option) {
	//发送锁
	sending := new(sync.Mutex)
	//处理请求等待
	wg := new(sync.WaitGroup)
	for {
		//读取请求
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		//处理请求
		go server.handleRequest(cc, req, sending, wg, opt.HandleTimeout)
	}
	wg.Wait()
	_ = cc.Close()
}

//请求结构体
type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
	mtype        *methodType
	svc          *service
}

//得到请求头
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

//查找服务
func (server *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	//得到服务名和方法名
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	//读取服务对象
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}

//得到请求数据结构体
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	//得到请求头
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	//根据服务名查找服务
	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}
	//新建请求对象
	req.argv = req.mtype.newArgv()
	//新建回复对象
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	//从conn读取请求数据到请求对象
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}
	return req, nil
}

//发送回复
func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	//往conn里面写入头和体
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

//处理请求
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	called := make(chan struct{})
	sent := make(chan struct{})
	go func() {
		//反射调用接口处理逻辑
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		//调用完成
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		//设置回复
		server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		//回复完成
		sent <- struct{}{}
	}()
	//没有超时，一直等待
	if timeout == 0 {
		<-called
		<-sent
		return
	}
	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", timeout)
		server.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}

//接收客户端连接
func (server *Server) Accept(lis net.Listener) {
	//循环接收
	for {
		//一个客户端连接
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		//处理交给协程
		go server.ServeConn(conn)
	}
}

//后端起服接收
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

//注册服务
func (server *Server) Register(rcvr interface{}) error {
	//新建一个rpc提供的服务
	s := newService(rcvr)
	//保存提供的服务
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

//注册服务
func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

const (
	connected        = "200 Connected to Gee RPC"
	//http转rpc请求路径
	defaultRPCPath   = "/_geeprc_"
	//普通http请求路径
	defaultDebugPath = "/debug/geerpc"
)

//用于rpc转换
func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//不允许非Connect方法
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = io.WriteString(w, "405 must CONNECT\n")
		return
	}
	//http连接转换为rpc连接
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}
	//回复前端200
	_, _ = io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	//处理tcp连接
	server.ServeConn(conn)
}

//开启http服务
func (server *Server) HandleHTTP() {
	http.Handle(defaultRPCPath, server)
	http.Handle(defaultDebugPath, debugHTTP{server})
}

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}
