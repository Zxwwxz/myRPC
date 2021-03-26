package httpBase

import (
    "context"
    "fmt"
    "github.com/gorilla/mux"
    "myRPC/config"
    "myRPC/log/base"
    "net"
    "net/http"
    "net/http/pprof"
    "strconv"
)

//http对象
type HttpServer struct {
    //tcp监听器
    listener  net.Listener
    //http服务器
    server    *http.Server
    //服务路径路由
    router    *mux.Router
}

func NewHttpServer(port int) (httpServer *HttpServer,err error) {
    listener,err := net.Listen("tcp",fmt.Sprintf(":%d",port))
    if err != nil {
        return nil,err
    }
    router := mux.NewRouter()
    httpServer = &HttpServer{listener:listener,router:router}
    return httpServer,nil
}

//开启
func (httpServer *HttpServer)Start()(err error)  {
    server := &http.Server{Handler:httpServer.router}
    httpServer.server = server
    //开启http
    return server.Serve(httpServer.listener)
}

//停止
func (httpServer *HttpServer)Stop()(err error)  {
    if httpServer.server != nil {
        return httpServer.server.Shutdown(context.TODO())
    }
    return nil
}

//获取路由，要在起服前就获取，进行添加路由
func (httpServer *HttpServer)GetRoute()(router *mux.Router)  {
    return httpServer.router
}

//默认开启的服务路由
func (httpServer *HttpServer)AddPropHandler()(err error)  {
    httpServer.router.HandleFunc("/debug/pprof/", pprof.Index)
    httpServer.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    httpServer.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
    httpServer.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    httpServer.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
    return nil
}

//监听参数修改
func (httpServer *HttpServer)AddParamsHandler(serviceConf *config.ServiceConf)(err error)  {
    httpServer.router.HandleFunc("/params", func(w http.ResponseWriter, r *http.Request){
        serviceHttpParams := config.ServiceHttpParams{}
        params := r.URL.Query()
        if params.Get("prometheus_switch_on") != ""{
            serviceHttpParams.PrometheusSwitchOn, _ = strconv.ParseBool(params.Get("prometheus_switch_on"))
        }
        if params.Get("client_limit_switch_on") != ""{
            serviceHttpParams.ClientLimitSwitchOn, _ = strconv.ParseBool(params.Get("limit_switch_on"))
        }
        if params.Get("server_limit_switch_on") != ""{
            serviceHttpParams.ServerLimitSwitchOn, _ = strconv.ParseBool(params.Get("server_limit_switch_on"))
        }
        if params.Get("trace_switch_on") != ""{
            serviceHttpParams.TraceSwitchOn, _ = strconv.ParseBool(params.Get("trace_switch_on"))
        }
        if params.Get("hystrix_switch_on") != "" {
            serviceHttpParams.HystrixSwitchOn, _ = strconv.ParseBool(params.Get("hystrix_switch_on"))
        }
        if params.Get("log_switch_on") != "" {
            serviceHttpParams.LogSwitchOn, _ = strconv.ParseBool(params.Get("log_switch_on"))
            _ = logBase.GetLogMgr().SetSwitchOn(serviceHttpParams.LogSwitchOn)
        }
        if params.Get("log_level") != "" {
            serviceHttpParams.LogLevel = params.Get("log_level")
            _ = logBase.GetLogMgr().SetLevel(serviceHttpParams.LogLevel)
        }
        serviceConf.HttpMergeConfig(serviceHttpParams)
        _,_ = w.Write([]byte("success"))
        fmt.Println("http serviceConf:",serviceConf)
    })
    return nil
}
