package generator

import (
	toolsBase "myRPC/template/base"
	"os"
	"path"
)

var dirList = []string{
	"client",
	"config",
	"controller",
	"proto",
	"docs",
	"logs",
	"model",
	"router",
	"scripts",
	"test",
	"test/config",
}

//目录生成器
type generatorDir struct {}

func NewGeneratorDir()(*generatorDir){
	return &generatorDir{}
}

func(g *generatorDir)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	for _,dir := range dirList{
		fullPath := path.Join(opt.OutputPath,dir)
		//分别创建文件
		err := os.MkdirAll(fullPath,0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func(g *generatorDir)Name() string{
	return "gen_dir"
}

