package netrpc

import (
	"bufio"
	"encoding/gob"
	"io"
)

func NewGoClientCodec(conn io.ReadWriteCloser) *GobClientCodec {
	encBuf := bufio.NewWriter(conn)
	return &GobClientCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(encBuf),
		encBuf: encBuf,
	}
}

type GobClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
}

//func (gcc *GobClientCodec) HookTrace(id uint64, reqMetaData string, f func(uint64, string)) interface{} {
//}

func (gcc *GobClientCodec) WriteRequest(r *Request, body interface{}) (err error) {
	if err = gcc.enc.Encode(r); err != nil {
		return
	}
	if err = gcc.enc.Encode(body); err != nil {
		return
	}
	return gcc.encBuf.Flush()
}

func (gcc *GobClientCodec) ReadResponseHeader(r *Response) error {
	return gcc.dec.Decode(r)
}

func (gcc *GobClientCodec) ReadResponseBody(body interface{}) error {
	return gcc.dec.Decode(body)
}

func (gcc *GobClientCodec) Close() error {
	return gcc.rwc.Close()
}
