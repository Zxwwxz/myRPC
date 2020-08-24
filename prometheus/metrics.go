package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/status"
)

var (
	defaultClientMetrics  *Metrics
	defaultServerMetrics  *Metrics
)

type Metrics struct {
	requestCounter    *prometheus.CounterVec
	errcodeCounter    *prometheus.CounterVec
	latencySummary    *prometheus.SummaryVec
}

func GetClientMetrics()(*Metrics)  {
	return defaultClientMetrics
}

func GetServerMetrics()(*Metrics)  {
	return defaultServerMetrics
}

func initMetrics()  {
	defaultClientMetrics = newServerMetrics()
	defaultServerMetrics = newServerMetrics()
}

//生成server metrics实例
func newServerMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "server_request",
			}, []string{"service", "method"}),
		errcodeCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "server_request_errcode",
			}, []string{"service", "method", "grpc_code"}),
		latencySummary: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "server_request_time",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			}, []string{"service", "method"},
		),
	}
}

//生成client metrics实例
func newClientMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_request",
			}, []string{"service", "method"}),
		errcodeCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_request_errcode",
			}, []string{"service", "method", "grpc_code"}),
		latencySummary: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "client_request_time",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			}, []string{"service", "method"},
		),
	}
}
func (m *Metrics) IncRequest(ctx context.Context, serviceName, methodName string) {
	m.requestCounter.WithLabelValues(serviceName, methodName).Inc()
}

func (m *Metrics) IncErrcode(ctx context.Context, serviceName, methodName string, err error) {
	st, _ := status.FromError(err)
	m.errcodeCounter.WithLabelValues(serviceName, methodName, st.Code().String()).Inc()
}

func (m *Metrics) ObserveLatency(ctx context.Context, serviceName, methodName string, useTime int64) {
	m.latencySummary.WithLabelValues(serviceName, methodName).Observe(float64(useTime))
}





