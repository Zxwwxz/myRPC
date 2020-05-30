package main

import "fmt"

type Options struct {
	Name string
	Age int
	Address string
	Money float32
}

func Init1(name string,age int,address string,money float32) Options {
	options := Options{}
	options.Name = name
	options.Age = age
	options.Address = address
	options.Money= money
	return options
}

func Init2(params ...interface{}) Options {
	options := Options{}
	for i,v := range params{
		switch i {
		case 0:
			value,ok := v.(string)
			if ok {options.Name = value}
		case 1:
			value,ok := v.(int)
			if ok {options.Age = value}
		case 2:
			value,ok := v.(string)
			if ok {options.Address = value}
		case 3:
			value,ok := v.(float32)
			if ok {options.Money = value}
		}
	}
	return options
}

type Option func(*Options)

func Init3(optionList ...Option) Options {
	options := Options{}
	for _,v := range optionList{
		v(&options)
	}
	return options
}

func SetName(name string) Option {
	return func(options *Options) {
		options.Name = name
	}
}

func SetAge(age int) Option {
	return func(options *Options) {
		options.Age = age
	}
}

func SetAddress(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func SetMoney(money float32) Option {
	return func(options *Options) {
		options.Money = money
	}
}

func main()  {
	fmt.Println("Init1:",Init1("aaa",10,"AAA",10.0))
	fmt.Println("Init2:",Init2("bbb",20,"BBB",20.0))
	fmt.Println("Init3:",Init3(SetName("ccc"),SetAge(30),SetAddress("CCC"),SetMoney(30.0)))
}
