package main

import (
	"fmt"
	"os"
	"text/template"
)

type Person struct {
	Name string
	Age  int
}

func main()  {
	t,err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Println("ParseFiles err：",err)
		return
	}
	person := &Person{Name:"zzz",Age:20}
	if err := t.Execute(os.Stdout, person); err != nil {
		fmt.Println("Execute err：", err)
	}
}
