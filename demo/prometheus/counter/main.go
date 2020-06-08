package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:"example_counter",
	})
)

func init()  {
	prometheus.MustRegister(counter)
}

func main()  {
	go func() {
		counter.Inc()
		time.Sleep(time.Second)
	}()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(":9100", nil))
}
