package main

import (
    "time"
    "net/http"
    _ "net/http/pprof"
)

func counterHttp() {
    slice := make([]int, 0)
    c := 1
    for i := 0; i < 100000; i++ {
        c = i + 1 + 2 + 3 + 4 + 5
        slice = append(slice, c)
    }
}

func workForever() {
    for {
        go counterHttp()
        time.Sleep(1 * time.Second)
    }
}

func httpGet(w http.ResponseWriter, r *http.Request) {
    counterHttp()
}

func main() {
    go workForever()
    http.HandleFunc("/get", httpGet)
    http.ListenAndServe(":9999", nil)
}
