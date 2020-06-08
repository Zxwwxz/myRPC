package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"time"
)

var (
	gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:"example_gauge",
	})
)

func init()  {
	prometheus.MustRegister(gauge)
}

func main()  {
	go func() {
		for{
			val := rand.Float64() * 100
			gauge.Set(val)
			time.Sleep(time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(":9100", nil))
}
