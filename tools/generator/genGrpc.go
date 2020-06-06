package generator

import (
	"fmt"
	toolsBase "myRPC/tools/base"
	"os/exec"
	"path"
)

type generatorGrpc struct {}

func NewGeneratorGrpc()(*generatorGrpc){
	return &generatorGrpc{}
}

func(g *generatorGrpc)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	fullPath := path.Join(opt.OutputPath,"generate/",opt.ProtoPath)
	_,err := toolsBase.CopyFile(opt.ProtoPath,fullPath)
	if err != nil {
		return err
	}
	fullParams := fmt.Sprintf("plugins=grpc:%s/generate/",opt.OutputPath)
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
