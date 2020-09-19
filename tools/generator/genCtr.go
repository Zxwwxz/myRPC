package generator

import (
	"github.com/ibinarytree/proto"
	"myRPC/tools/base"
	"myRPC/util"
	"os"
	"path"
	"text/template"
)

type generatorCtr struct {}

func NewGeneratorCtr()(*generatorCtr){
	return &generatorCtr{}
}

func(g *generatorCtr)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	filename := path.Join(opt.OutputPath, "controller/controller.go")
	if util.IsFileExist(filename) {
		return nil
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	t := template.New("Ctr")
	t, err = t.Parse(ctrTemplateFile)
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
		filename := path.Join(opt.OutputPath, "controller/"+rpc.Name+".go" )
		if util.IsFileExist(filename) {
			continue
		}
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			continue
		}
		defer file.Close()
		t := template.New("Ctr"+rpc.Name)
		tempFile := ctrTemplateFuncFile
		if rpc.StreamsRequest == true && rpc.StreamsReturns == true {
			tempFile = ctrTemplateStreamFuncFile
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

func(g *generatorCtr)Name() string{
	return "gen_Ctr"
}
