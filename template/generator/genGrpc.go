package generator

import (
	"fmt"
	toolsBase "myRPC/template/base"
	"os/exec"
	"path"
)

type generatorGrpc struct {}

func NewGeneratorGrpc()(*generatorGrpc){
	return &generatorGrpc{}
}

func(g *generatorGrpc)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	fullPath := path.Join(opt.OutputPath,"proto/",opt.ProtoPath)
	_,err := toolsBase.CopyFile(opt.ProtoPath,fullPath)
	if err != nil {
		return err
	}
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
