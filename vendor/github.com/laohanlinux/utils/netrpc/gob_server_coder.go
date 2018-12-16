package netrpc

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
	closed bool

	// extra attributes
	reqMetaData   interface{}
	replyMetaData interface{}
}

func (gsc *GobServerCodec) SetReqMetaData(metaData interface{}) {
	gsc.reqMetaData = metaData
}

func (gsc *GobServerCodec) GetReplyMetaData() interface{} {
	return gsc.replyMetaData
}

func (gsc *GobServerCodec) ReadRequestHeader(r *Request) error {
	return gsc.dec.Decode(r)
}

func (gsc *GobServerCodec) ReadRequestBody(body interface{}) error {
	return gsc.dec.Decode(body)
}

func (gsc *GobServerCodec) WriteResponse(r *Response, body interface{}) (err error) {
	if err = gsc.enc.Encode(r); err != nil {
		if gsc.encBuf.Flush() == nil {
			log.Printf("rpc gob error encoding:%v\n", err)
			gsc.Close()
		}
		return
	}
	if err = gsc.enc.Encode(body); err != nil {
		if gsc.encBuf.Flush() == nil {
			log.Printf("rpc gob error encoding body:%v\n", err)
			gsc.Close()
		}
		return
	}
	return gsc.encBuf.Flush()
}

func (gsc *GobServerCodec) Close() error {
	if gsc.closed {
		return nil
	}
	gsc.closed = true
	return gsc.rwc.Close()
}
