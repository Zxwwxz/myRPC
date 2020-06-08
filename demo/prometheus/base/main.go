package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main()  {
	flag.Parse()
	http.Handle("/metrics",promhttp.Handler())
	fmt.Println(http.ListenAndServe(":9100", nil))
}
