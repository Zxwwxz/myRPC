package generator

import (
	"myRPC/template/base"
	"os"
	"path"
	"text/template"
)

type generatorConfig struct {}

func NewGeneratorConfig()(*generatorConfig){
	return &generatorConfig{}
}

func(g *generatorConfig)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	err := g.CreateFile(opt,meta,"config/")
	err = g.CreateFile(opt,meta,"test/config/")
	return err
}

func(g *generatorConfig)CreateFile(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData,dir string)(err error){
	filename := path.Join(opt.OutputPath, dir, "/config.yaml")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	t := template.New("Config")
	t, err = t.Parse(configTemplateFile)
	if err != nil {
		return err
	}
	err = t.Execute(file, meta)
	if err != nil {
		return err
	}
	return nil
}

func(g *generatorConfig)Name() string{
	return "gen_config"
}
