package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

const (
	default_start = 100
	default_width = 100
	default_count = 5
)

type Metrics struct {
	//请求数
	requestCounter    *prometheus.CounterVec
	//错误数
	errcodeCounter    *prometheus.CounterVec
	//请求耗时
	latencyHistogram  *prometheus.HistogramVec
}

func NewMetrics(clientHistogram,serverHistogram string)(*Metrics,*Metrics)  {
	defaultClientMetrics := newClientMetrics(getHistogram(clientHistogram))
	defaultServerMetrics := newServerMetrics(getHistogram(serverHistogram))
	return defaultClientMetrics,defaultServerMetrics
}

//请求耗时的采样率参数
func getHistogram(histogram string)(start,width float64,count int)  {
	start = default_start
	width = default_width
	count = default_count
	histogramSlice := strings.Split(histogram,",")
	if len(histogramSlice) != 3 {
		return start,width,count
	}
	if temp, err := strconv.ParseFloat(histogramSlice[0], 64);err == nil {
		start = temp
	}
	if temp, err := strconv.ParseFloat(histogramSlice[1], 64);err == nil {
		width = temp
	}
	if temp, err := strconv.Atoi(histogramSlice[2]) ;err == nil {
		count = temp
	}
	return start,width,count
}

//生成server metrics实例
func newServerMetrics(start,width float64,count int) *Metrics {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_request",
		}, []string{"service", "method"})
	errcodeCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_request_errcode",
		}, []string{"service", "method", "type", "grpc_code"})
	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:       "server_request_time",
			Help:       "RPC latency distributions.",
			Buckets:    prometheus.LinearBuckets(start,width,count),
		}, []string{"service", "method"})
	//注册指标
	prometheus.MustRegister(requestCounter,errcodeCounter,latencyHistogram)
	return &Metrics{
		requestCounter: requestCounter,
		errcodeCounter: errcodeCounter,
		latencyHistogram: latencyHistogram,
	}
}

//生成client metrics实例
func newClientMetrics(start,width float64,count int) *Metrics {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "client_request",
		}, []string{"service", "method"})
	errcodeCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "client_request_errcode",
		}, []string{"service", "method", "grpc_code"})
	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:       "client_request_time",
			Help:       "RPC latency distributions.",
			Buckets:    prometheus.LinearBuckets(start,width,count),
		}, []string{"service", "method"})
	prometheus.MustRegister(requestCounter,errcodeCounter,latencyHistogram)
	return &Metrics{
		requestCounter: requestCounter,
		errcodeCounter: errcodeCounter,
		latencyHistogram: latencyHistogram,
	}
}

//增加某服务，某接口的请求数
func (m *Metrics) IncRequest(ctx context.Context, serviceName, methodName string) {
	m.requestCounter.WithLabelValues(serviceName, methodName).Inc()
}

//增加某服务，某接口，服务端还是客户端，哪个错误码的错误数
func (m *Metrics) IncErrcode(ctx context.Context, serviceName, methodName string, incType string, code int) {
	m.errcodeCounter.WithLabelValues(serviceName, methodName, incType, strconv.Itoa(code)).Inc()
}

//增加某服务，某接口的调用时间
func (m *Metrics) ObserveLatency(ctx context.Context, serviceName, methodName string, useTime int64) {
	m.latencyHistogram.WithLabelValues(serviceName, methodName).Observe(float64(useTime))
}