package toolsBase

import "github.com/ibinarytree/proto"

//模板参数
type Option struct {
	//pb路径：xxx.proto
	ProtoPath      string
	//服务输出路径：../../myRPC/myService/xxx
	OutputPath     string
}

//解析pb文件后的参数
type ServiceMetaData struct {
	//pb中的服务对象 service xxx
	Service  *proto.Service
	//目前没用到
	Messages []*proto.Message
	//pb中的接口对象列表 rpc xxx
	Rpc      []*proto.RPC
	//pb中的包对象 package xxx
	Package  *proto.Package
	//服务名，对应pb文件名 xxx
	ServiceName      string
	//服务代码导入前置路径：myRPC/myService
	ImportPreFix     string
}

