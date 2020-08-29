package generator

import (
	"github.com/ibinarytree/proto"
	"myRPC/tools/base"
	"myRPC/util"
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
	for _,rpc := range meta.Rpc{
		newMeta := *meta
		newMeta.Rpc = []*proto.RPC{rpc}
		filename := path.Join(opt.OutputPath, "client/"+rpc.Name+".go" )
		if util.IsFileExist(filename) {
			continue
		}
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			continue
		}
		defer file.Close()
		t := template.New("client"+rpc.Name)
		t, err = t.Parse(clientTemplateFuncFile)
		if err != nil {
			continue
		}
		err = t.Execute(file, newMeta)
		if err != nil {
			continue
		}
	}
	return nil
}

func(g *generatorClient)Name() string{
	return "gen_client"
}
