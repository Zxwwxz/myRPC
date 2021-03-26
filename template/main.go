package main

import (
	"fmt"
	"github.com/urfave/cli"
	toolsBase "myRPC/template/base"
	"myRPC/template/generator"
	"os"
)

//作用：用于生成服务模板
//定义：将要生成的服务，定义成pb文件，放到当前同级目录下
//命令：go run main.go -p xxx.proto(pb路径) -o path/xxx(服务输出路径，非必须)
func main()  {
	var opt toolsBase.Option
	var meta toolsBase.ServiceMetaData
	//参数解析器
	app := cli.NewApp()
	app.Version = "2.0.1"
	//接收两个参数 -o和-p
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "p",
			Required:    true,
			Destination: &opt.ProtoPath,
		},
		cli.StringFlag{
			Name:        "o",
			Destination: &opt.OutputPath,
		},
	}
	app.Action = func(c *cli.Context) error {
		//注册生成器
		generator.GeneratorMgrObj.Register()
		//执行生成器
		generator.GeneratorMgrObj.Run(&opt,&meta)
		return nil
	}
	//开始执行
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("app Run err:",err)
	}
}
