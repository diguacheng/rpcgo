package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// 这里主要是具体编码的实现 实现 Codec 

type GobCodec struct {
	conn io.ReadWriteCloser  // 由构建函数传入通常是通过 TCP 或者 Unix 建立 socket 时得到的链接实例
	buf *bufio.Writer  //buf 是为了防止阻塞而创建的带缓冲的 Writer，一般这么做能提升性能。
	dec *gob.Decoder
	enc *gob.Encoder
}

var _ Codec =(*GobCodec)(nil) // 这里主要的作用是 检查 GobCodec 是否实现了Godec的接口 

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *GobCodec) ReadHeader(h *Header) error { //将数据解码 
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush() // 将缓存内的数据 写入conn 中 
		if err != nil {
			_ = c.Close()
		}
	}()
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}