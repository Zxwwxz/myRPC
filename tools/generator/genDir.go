package generator

import (
	toolsBase "myRPC/tools/base"
	"os"
	"path"
)

var dirList = []string{
	"controller",
	"ild",
	"scripts",
	"conf",
	"app/router",
	"app/config",
	"model",
	"generate",
	"main",
}

type generatorDir struct {}

func NewGeneratorDir()(*generatorDir){
	return &generatorDir{}
}

func(g *generatorDir)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	for _,dir := range dirList{
		fullPath := path.Join(opt.OutputPath,dir)
		err := os.MkdirAll(fullPath,0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func(g *generatorDir)Name() string{
	return "gen_dir"
}

