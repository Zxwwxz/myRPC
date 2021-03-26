package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

//Gob编解码对象
type GobCodec struct {
	//连接对象
	conn io.ReadWriteCloser
	//发送时候的缓冲通道
	buf  *bufio.Writer
	//解码对象
	dec  *gob.Decoder
	//编码对象
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

//新建Gob编解码对象，需要依赖连接对象
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

//客户端服务端编码解码完全对称
//读取包头包体，只要编码的时候是用gob编码的，就能解析出来
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

//写入包头包体
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		//写入完成要刷新缓冲区
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

//关闭，直接关闭连接对象
func (c *GobCodec) Close() error {
	return c.conn.Close()
}
