// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat.proto

package pbChat

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ReqStruct struct {
	ReqMsg               string   `protobuf:"bytes,1,opt,name=reqMsg,proto3" json:"reqMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqStruct) Reset()         { *m = ReqStruct{} }
func (m *ReqStruct) String() string { return proto.CompactTextString(m) }
func (*ReqStruct) ProtoMessage()    {}
func (*ReqStruct) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{0}
}

func (m *ReqStruct) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqStruct.Unmarshal(m, b)
}
func (m *ReqStruct) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqStruct.Marshal(b, m, deterministic)
}
func (m *ReqStruct) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqStruct.Merge(m, src)
}
func (m *ReqStruct) XXX_Size() int {
	return xxx_messageInfo_ReqStruct.Size(m)
}
func (m *ReqStruct) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqStruct.DiscardUnknown(m)
}

var xxx_messageInfo_ReqStruct proto.InternalMessageInfo

func (m *ReqStruct) GetReqMsg() string {
	if m != nil {
		return m.ReqMsg
	}
	return ""
}

type RspStruct struct {
	RspMsg               string   `protobuf:"bytes,1,opt,name=rspMsg,proto3" json:"rspMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RspStruct) Reset()         { *m = RspStruct{} }
func (m *RspStruct) String() string { return proto.CompactTextString(m) }
func (*RspStruct) ProtoMessage()    {}
func (*RspStruct) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{1}
}

func (m *RspStruct) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RspStruct.Unmarshal(m, b)
}
func (m *RspStruct) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RspStruct.Marshal(b, m, deterministic)
}
func (m *RspStruct) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RspStruct.Merge(m, src)
}
func (m *RspStruct) XXX_Size() int {
	return xxx_messageInfo_RspStruct.Size(m)
}
func (m *RspStruct) XXX_DiscardUnknown() {
	xxx_messageInfo_RspStruct.DiscardUnknown(m)
}

var xxx_messageInfo_RspStruct proto.InternalMessageInfo

func (m *RspStruct) GetRspMsg() string {
	if m != nil {
		return m.RspMsg
	}
	return ""
}

func init() {
	proto.RegisterType((*ReqStruct)(nil), "pbChat.ReqStruct")
	proto.RegisterType((*RspStruct)(nil), "pbChat.RspStruct")
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor_8c585a45e2093e54) }

var fileDescriptor_8c585a45e2093e54 = []byte{
	// 130 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0xce, 0x48, 0x2c,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x48, 0x72, 0xce, 0x48, 0x2c, 0x51, 0x52,
	0xe6, 0xe2, 0x0c, 0x4a, 0x2d, 0x0c, 0x2e, 0x29, 0x2a, 0x4d, 0x2e, 0x11, 0x12, 0xe3, 0x62, 0x2b,
	0x4a, 0x2d, 0xf4, 0x2d, 0x4e, 0x97, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0xf2, 0xc0, 0x8a,
	0x8a, 0x0b, 0x90, 0x14, 0x15, 0x17, 0x20, 0x2b, 0x02, 0xf3, 0x8c, 0x1c, 0xb8, 0xb8, 0x40, 0x26,
	0x06, 0xa7, 0x16, 0x95, 0xa5, 0x16, 0x09, 0x19, 0x71, 0xb1, 0x80, 0x78, 0x42, 0x82, 0x7a, 0x10,
	0x8b, 0xf4, 0xe0, 0xb6, 0x48, 0x21, 0x84, 0x60, 0x66, 0x2a, 0x31, 0x68, 0x30, 0x1a, 0x30, 0x26,
	0xb1, 0x81, 0x9d, 0x66, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x80, 0xcb, 0x19, 0x9a, 0xa8, 0x00,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChatServerClient is the client API for ChatServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatServerClient interface {
	Chat(ctx context.Context, opts ...grpc.CallOption) (ChatServer_ChatClient, error)
}

type chatServerClient struct {
	cc *grpc.ClientConn
}

func NewChatServerClient(cc *grpc.ClientConn) ChatServerClient {
	return &chatServerClient{cc}
}

func (c *chatServerClient) Chat(ctx context.Context, opts ...grpc.CallOption) (ChatServer_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ChatServer_serviceDesc.Streams[0], "/pbChat.ChatServer/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServerChatClient{stream}
	return x, nil
}

type ChatServer_ChatClient interface {
	Send(*ReqStruct) error
	Recv() (*RspStruct, error)
	grpc.ClientStream
}

type chatServerChatClient struct {
	grpc.ClientStream
}

func (x *chatServerChatClient) Send(m *ReqStruct) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServerChatClient) Recv() (*RspStruct, error) {
	m := new(RspStruct)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServerServer is the server API for ChatServer service.
type ChatServerServer interface {
	Chat(ChatServer_ChatServer) error
}

func RegisterChatServerServer(s *grpc.Server, srv ChatServerServer) {
	s.RegisterService(&_ChatServer_serviceDesc, srv)
}

func _ChatServer_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServerServer).Chat(&chatServerChatServer{stream})
}

type ChatServer_ChatServer interface {
	Send(*RspStruct) error
	Recv() (*ReqStruct, error)
	grpc.ServerStream
}

type chatServerChatServer struct {
	grpc.ServerStream
}

func (x *chatServerChatServer) Send(m *RspStruct) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServerChatServer) Recv() (*ReqStruct, error) {
	m := new(ReqStruct)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ChatServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pbChat.ChatServer",
	HandlerType: (*ChatServerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _ChatServer_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "chat.proto",
}
