package generator

import (
	"github.com/ibinarytree/proto"
	"myRPC/template/base"
	"myRPC/util"
	"os"
	"path"
	"text/template"
)

//根据目标生成所有client文件
type generatorClient struct {}

func NewGeneratorClient()(*generatorClient){
	return &generatorClient{}
}

func(g *generatorClient)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	//创建所有接口共有的文件
	filename := path.Join(opt.OutputPath, "client/client.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
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
	//每个接口创建一个文件
	for _,rpc := range meta.Rpc{
		newMeta := *meta
		newMeta.Rpc = []*proto.RPC{rpc}
		filename := path.Join(opt.OutputPath, "client/"+rpc.Name+".go" )
		//已经存在不再创建，防止覆盖写好的代码，新增接口时才生成
		if util.IsFileExist(filename) {
			continue
		}
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			continue
		}
		defer file.Close()
		t := template.New("client"+rpc.Name)
		//模板需要区分普通模式和流模式
		tempFile := clientTemplateFuncFile
		if rpc.StreamsRequest == true && rpc.StreamsReturns == true {
			tempFile = clientTemplateStreamFuncFile
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

func(g *generatorClient)Name() string{
	return "gen_client"
}
