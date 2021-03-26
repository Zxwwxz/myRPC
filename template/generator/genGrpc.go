package generator

import (
	"fmt"
	toolsBase "myRPC/template/base"
	"myRPC/util"
	"os/exec"
	"path"
)

//pb文件拷贝并编译成go
type generatorGrpc struct {}

func NewGeneratorGrpc()(*generatorGrpc){
	return &generatorGrpc{}
}

func(g *generatorGrpc)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	fullPath := path.Join(opt.OutputPath,"proto/",opt.ProtoPath)
	//拷贝pb文件到服务目录下
	_,err := util.CopyFile(opt.ProtoPath,fullPath)
	if err != nil {
		return err
	}
	//生成pb对应的go文件
	fullParams := fmt.Sprintf("plugins=grpc:%s/proto/",opt.OutputPath)
	cmd := exec.Command("protoc","--go_out",fullParams,opt.ProtoPath)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func(g *generatorGrpc)Name() string{
	return "gen_grpc"
}
