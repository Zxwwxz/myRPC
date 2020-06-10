package service

import (
	"fmt"
	"myRPC/config"
)

func Init()  {
	err := config.InitConfig()
	if err != nil {
		fmt.Println("InitConfig,err:",err)
	}
}
