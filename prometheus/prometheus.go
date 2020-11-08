package prometheus

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"myRPC/config"
	"net/http"
)

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
