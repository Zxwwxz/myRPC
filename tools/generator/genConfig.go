package generator

import (
	"myRPC/util"
	toolsBase "myRPC/tools/base"
	"os"
	"path"
	"text/template"
)

type generatorConfig struct {}

func NewGeneratorConfig()(*generatorConfig){
	return &generatorConfig{}
}

func(g *generatorConfig)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	envStr := util.GetEnv()
	err := g.CreateFile(envStr,opt,meta)
	return err
}

func(g *generatorConfig)CreateFile(env string,opt *toolsBase.Option,meta *toolsBase.ServiceMetaData)(err error){
	filename := path.Join(opt.OutputPath, "config/",env, "/config.yaml")
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
