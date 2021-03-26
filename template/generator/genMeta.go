package generator

import (
	"github.com/ibinarytree/proto"
	toolsBase "myRPC/template/base"
	"os"
	"path"
)

//解析pb参数生成器
type generatorMeta struct {
	meta *toolsBase.ServiceMetaData
}

func NewGeneratorMeta()(*generatorMeta){
	return &generatorMeta{}
}

func(g *generatorMeta)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	g.meta = meta
	//打开pb文件
	reader, err := os.Open(opt.ProtoPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	//解析pb文件
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

//得到服务对象
func (g *generatorMeta) handleService(s *proto.Service) {
	g.meta.Service = s
}

//得到每个msg对象
func (g *generatorMeta) handleMessage(m *proto.Message) {
	g.meta.Messages = append(g.meta.Messages, m)
}

//得到每个rpc对象
func (g *generatorMeta) handleRPC(r *proto.RPC) {
	g.meta.Rpc = append(g.meta.Rpc, r)
}

//得到包对象
func (g *generatorMeta) handlePackage(r *proto.Package) {
	g.meta.Package = r
}

func (g *generatorMeta) handleOtherParams(opt *toolsBase.Option) {
	serviceName := opt.ProtoPath[:len(opt.ProtoPath)-6]
	//服务名就是pb文件名
	g.meta.ServiceName = serviceName
	if opt.OutputPath == ""{
		opt.OutputPath = toolsBase.OutPath
	}
	g.meta.ImportPreFix = path.Join(opt.OutputPath,serviceName)
	opt.OutputPath = path.Join("../../",opt.OutputPath,serviceName)
}
