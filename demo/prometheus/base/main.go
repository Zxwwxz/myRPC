package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main()  {
	flag.Parse()
	//普罗米修斯平台拉取当前服务的:9100/metrics，会获取到本服务器的性能数据，进行可视化展示
	http.Handle("/metrics",promhttp.Handler())
	fmt.Println(http.ListenAndServe(":9100", nil))
}
