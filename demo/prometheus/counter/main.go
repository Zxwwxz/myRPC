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
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:"example_counter",
	})
	rand.Seed(time.Now().Unix())
	prometheus.MustRegister(counter)
	go func() {
		for{
			val := rand.Float64()
			counter.Add(val)
			time.Sleep(time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(":9100", nil))
}
