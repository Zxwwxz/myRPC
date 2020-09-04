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
    histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:"example_Histogram",
        Buckets:prometheus.LinearBuckets(0.2,0.2,5),
    },[]string{"serName"})
    rand.Seed(time.Now().Unix())
    prometheus.MustRegister(histogram)
    go func() {
        for{
            val := rand.Float64()
            histogram.WithLabelValues("serName1").Observe(val)
            time.Sleep(time.Second)
        }
    }()
    http.Handle("/metrics", promhttp.Handler())
    fmt.Println(http.ListenAndServe(":9100", nil))
}
