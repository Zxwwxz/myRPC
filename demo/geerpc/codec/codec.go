package codec

import (
	"io"
)

//| Option | Header1 | Body1 | Header2 | Body2 |
//rpc请求的头部，每次请求都会携带
type Header struct {
	//请求服务和接口
	ServiceMethod string
	//客户端请求序号
	Seq           uint64
	//服务端是否有报错
	Error         string
}

type Codec interface {
	//连接对象
	io.Closer
	//读取包头，将空的头部传入，里面会进行解析，外面再去拿到的头部就会有值
	ReadHeader(*Header) error
	//读取包体，将空的对象传入，里面会进行解析，外面再去拿到的对象就会有值
	ReadBody(interface{}) error
	//写入包头和包体，将有值的包头和包体传入，里面会直接发送到连接对象中
	Write(*Header, interface{}) error
}

//编解码函数，调用函数，连接对象当做参数，会返回得到一个编解码对象
type NewCodecFunc func(io.ReadWriteCloser) Codec

//编解码类型
type Type string

const (
	//编解码类型-Gob
	GobType  Type = "application/gob"
	//编解码类型-Json
	JsonType Type = "application/json"
)

//编解码函数map，根据编解码类型，得到一个编解码函数
//这里为什么不直接根据编解码类型，得到一个编解码对象？
//因为编解码对象需要依赖连接对象conn，所以要返回一个函数，调用函数传入conn，得到对象
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	//初始化编解码函数map
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	//这里只添加了Gob类型的编解码对象，如果要支持Json或Pb只需要改这里加类型，同时加一个Codec接口实现
	NewCodecFuncMap[GobType] = NewGobCodec
}
