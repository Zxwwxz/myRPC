package toolsBase

import "github.com/ibinarytree/proto"

type Option struct {
	ProtoPath      string
	OutputPath     string
	GenClientCode  bool
	GenServerCode  bool
}

type ServiceMetaData struct {
	Service  *proto.Service
	Messages []*proto.Message
	Rpc      []*proto.RPC
	Package  *proto.Package
	ServiceName      string
}

