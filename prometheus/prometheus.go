package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func InitPrometheus(listenPort int,clientHistogram,serverHistogram string)(error) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ =  http.ListenAndServe(fmt.Sprint(":%d",listenPort), nil)
	}()
	initMetrics(clientHistogram,serverHistogram)
	return nil
}
