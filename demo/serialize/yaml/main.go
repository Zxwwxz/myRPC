package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ClientConfig Client     `yaml:"client"`
	ServiceConfig Service   `yaml:"service"`
}

type Client struct {
	Ip []string     `yaml:"ip"`
	Port int        `yaml:"port"`
}

type Service struct {
	Type int      `yaml:"type"`
	Id int        `yaml:"id"`
	Name string   `yaml:"name"`
}

func main()  {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Printf("unmarshal failed err:%v\n", err)
		return
	}

	fmt.Printf("site port:%d\n", conf.ClientConfig.Ip)
}