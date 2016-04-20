package msgpack

import (
	"encoding/json"
	"time"
)

type ResponseMsg struct {
	Results interface{} `json:"results, omitempty"`
	ErrCode int         `json:"error, omitempty"`
	Time    float64     `json:"time,omitempty"`
	start   time.Time
}

func (msg *ResponseMsg) JsonToBytes(errCode ...int) []byte {
	if len(errCode) > 0 {
		msg.ErrCode = errCode[0]
	}
	b, _ := json.Marshal(msg)
	return b
}
