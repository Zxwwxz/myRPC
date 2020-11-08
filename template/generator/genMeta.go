package generator

import (
	"github.com/ibinarytree/proto"
	toolsBase "myRPC/template/base"
	"os"
	"path"
)

type generatorMeta struct {
	meta *toolsBase.ServiceMetaData
}

func NewGeneratorMeta()(*generatorMeta){
	return &generatorMeta{}
}

func(g *generatorMeta)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	g.meta = meta
	reader, err := os.Open(opt.ProtoPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return err
	}
	proto.Walk(definition,
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC),
		proto.WithPackage(g.handlePackage),
	)
	g.handleOtherParams(opt)
	return nil
}

func(g *generatorMeta)Name() string{
	return "gen_meta"
}

func (g *generatorMeta) handleService(s *proto.Service) {
	g.meta.Service = s
}

func (g *generatorMeta) handleMessage(m *proto.Message) {
	g.meta.Messages = append(g.meta.Messages, m)
}

func (g *generatorMeta) handleRPC(r *proto.RPC) {
	g.meta.Rpc = append(g.meta.Rpc, r)
}

func (g *generatorMeta) handlePackage(r *proto.Package) {
	g.meta.Package = r
}

func (g *generatorMeta) handleOtherParams(opt *toolsBase.Option) {
	serviceName := opt.ProtoPath[:len(opt.ProtoPath)-6]
	g.meta.ServiceName = serviceName
	if opt.OutputPath == ""{
		opt.OutputPath = toolsBase.OutPath
	}
	opt.ImportPreFix = toolsBase.OutPath
	opt.OutputPath = path.Join("../../",opt.OutputPath,serviceName)
	g.meta.ImportPreFix = path.Join(opt.ImportPreFix,serviceName)
}
