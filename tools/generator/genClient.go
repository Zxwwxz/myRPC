package generator

import (
	toolsBase "myRPC/tools/base"
	"os"
	"path"
	"text/template"
)

type generatorClient struct {}

func NewGeneratorClient()(*generatorClient){
	return &generatorClient{}
}

func(g *generatorClient)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	filename := path.Join(opt.OutputPath, "client/client.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	t := template.New("client")
	t, err = t.Parse(clientTemplateFile)
	if err != nil {
		return err
	}
	err = t.Execute(file, meta)
	if err != nil {
		return err
	}
	return nil
}

func(g *generatorClient)Name() string{
	return "gen_client"
}
