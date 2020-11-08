package main

import (
    "fmt"
    "reflect"
)

type AStruct struct {
    B BStruct
}

type BStruct struct {
    C int
}

func main() {
    aaa := 100
    a := AStruct{
        B:BStruct{C:aaa},
    }
    reflect.ValueOf(a).FieldByName("B").FieldByName("C").Elem().Set(reflect.ValueOf(80))
    fmt.Println("a:",a)
}
