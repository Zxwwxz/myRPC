package prometheus

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"myRPC/config"
	"net/http"
)

//当前监控数据是每次调用会进行修改
//服务器开启普罗米修斯服务后，会定时到当前服务拉取数据，进行统计分析
//grafana可以设置数据来源是普罗米修斯，进行统计查看
var prometheusManager *PrometheusManager

type PrometheusManager struct {
	defaultClientMetrics  *Metrics
	defaultServerMetrics  *Metrics
}

func NewPrometheusManager(clientHistogram,serverHistogram string)(error) {
	prometheusManager = &PrometheusManager{}
	prometheusManager.defaultClientMetrics,prometheusManager.defaultServerMetrics = NewMetrics(clientHistogram,serverHistogram)
	return nil
}

func GetPrometheusManager()(*PrometheusManager) {
	return prometheusManager
}

func (prometheusManager *PrometheusManager)GetClientMetrics()(*Metrics)  {
	return prometheusManager.defaultClientMetrics
}

func (prometheusManager *PrometheusManager)GetServerMetrics()(*Metrics)  {
	return prometheusManager.defaultServerMetrics
}

//开启http端口监听
func (prometheusManager *PrometheusManager)AddPrometheusHandler(router *mux.Router,serviceConf *config.ServiceConf)(err error)  {
	router.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		if serviceConf.Prometheus.SwitchOn{
			promhttp.Handler().ServeHTTP(writer,request)
		}
	})
	return nil
}

func (prometheusManager *PrometheusManager)Stop()  {
}
