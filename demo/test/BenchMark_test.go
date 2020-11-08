package test

import (
    "fmt"
    "strconv"
    "testing"
)

//go test -bench=. -benchmem -run=none
//函数名-GOMAXPROCS数       执行次数        每次执行时间 ns/op        每次分配字节数 B/op        每次分配内存数 allocs/op
func BenchmarkSprintf(b *testing.B){
    num:=10
    b.ResetTimer()
    for i:=0;i<b.N;i++{
        fmt.Sprintf("%d",num)
    }
}

func BenchmarkFormat(b *testing.B){
    num:=int64(10)
    b.ResetTimer()
    for i:=0;i<b.N;i++{
        strconv.FormatInt(num,10)
    }
}

func BenchmarkItoa(b *testing.B){
    num:=10
    b.ResetTimer()
    for i:=0;i<b.N;i++{
        strconv.Itoa(num)
    }
}
