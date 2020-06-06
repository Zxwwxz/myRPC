package generator

import (
	toolsBase "myRPC/tools/base"
	"os"
	"path"
	"text/template"
)

type generatorRouter struct {}

func NewGeneratorRouter()(*generatorRouter){
	return &generatorRouter{}
}

func(g *generatorRouter)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	filename := path.Join(opt.OutputPath, "router/router.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	t := template.New("Router")
	t, err = t.Parse(routerTemplateFile)
	if err != nil {
		return err
	}
	err = t.Execute(file, meta)
	if err != nil {
		return err
	}
	return nil
}

func(g *generatorRouter)Name() string{
	return "gen_Router"
}
