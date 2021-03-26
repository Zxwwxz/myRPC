package generator

import (
	"github.com/ibinarytree/proto"
	toolsBase "myRPC/template/base"
	"myRPC/util"
	"os"
	"path"
	"text/template"
)

//根据目标生成所有router文件
type generatorRouter struct {}

func NewGeneratorRouter()(*generatorRouter){
	return &generatorRouter{}
}

func(g *generatorRouter)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	//创建所有接口共有的文件
	filename := path.Join(opt.OutputPath, "router/router.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
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
	//每个接口创建一个文件
	for _,rpc := range meta.Rpc {
		newMeta := *meta
		newMeta.Rpc = []*proto.RPC{rpc}
		filename := path.Join(opt.OutputPath, "router/"+rpc.Name+".go" )
		if util.IsFileExist(filename) {
			continue
		}
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		//已经存在不再创建，防止覆盖写好的代码，新增接口时才生成
		if err != nil {
			continue
		}
		defer file.Close()
		t := template.New("Router"+rpc.Name)
		//模板需要区分普通模式和流模式
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
