package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func InitPrometheus(listenPort int)(error) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ =  http.ListenAndServe(fmt.Sprint(":%d",listenPort), nil)
	}()
	initMetrics()
	return nil
}
