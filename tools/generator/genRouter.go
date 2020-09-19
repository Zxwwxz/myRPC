package generator

import (
	"github.com/ibinarytree/proto"
	toolsBase "myRPC/tools/base"
	"myRPC/util"
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
	for _,rpc := range meta.Rpc {
		newMeta := *meta
		newMeta.Rpc = []*proto.RPC{rpc}
		filename := path.Join(opt.OutputPath, "router/"+rpc.Name+".go" )
		if util.IsFileExist(filename) {
			continue
		}
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			continue
		}
		defer file.Close()

		t := template.New("Router"+rpc.Name)
		tempFile := routerTemplateFuncFile
		if rpc.StreamsRequest == true && rpc.StreamsReturns == true {
			tempFile = routerTemplateStreamFuncFile
		}
		t, err = t.Parse(tempFile)
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

func(g *generatorRouter)Name() string{
	return "gen_Router"
}
