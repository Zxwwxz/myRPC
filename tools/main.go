package main

import (
	"fmt"
	"github.com/urfave/cli"
	toolsBase "myRPC/tools/base"
	"myRPC/tools/generator"
	"os"
)

func main()  {
	var opt toolsBase.Option
	var meta toolsBase.ServiceMetaData
	app := cli.NewApp()
	app.Version = "2.0.1"

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
		cli.StringFlag{
			Name:        "i",
			Destination: &opt.ImportPreFix,
		},
		cli.BoolFlag{
			Name:        "c",
			Destination: &opt.GenClientCode,
		},
		cli.BoolFlag{
			Name:        "s",
			Destination: &opt.GenServerCode,
		},
	}
	app.Action = func(c *cli.Context) error {
		generator.GeneratorMgrObj.Register()
		generator.GeneratorMgrObj.Run(&opt,&meta)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("app Run err:",err)
	}
}
