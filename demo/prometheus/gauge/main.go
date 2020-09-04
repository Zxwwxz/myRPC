package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"time"
)

func main()  {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:"example_gauge",
	})
	rand.Seed(time.Now().Unix())
	prometheus.MustRegister(gauge)
	go func() {
		for{
			val := rand.Float64()
			gauge.Set(val)
			time.Sleep(time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(":9100", nil))
}
