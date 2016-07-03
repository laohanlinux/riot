package handler

import "encoding/json"

const (
	Ok             = "ok"
	OpErr          = "operation errror"
	NetErr         = "net work timeout"
	InternalErr    = "riot server error"
	InvalidKey     = "invalid key"
	InvalidRequest = "invalid request"
	QsInvalid      = "invalid query strategies"
	NotFound       = "not found"
	NotExistBucket = "not exist the bucket"
	InvalidBucket  = "invalid bucket"
)

var MsgErrCodeMap map[string]ErrCodeObj

type ErrCodeObj struct {
	httpCode   int
	StatusCode int    `json:"errCode"`
	Info       string `json:"msg"`
}

func (err *ErrCodeObj) toJSONBytes() []byte {
	b, _ := json.Marshal(err)
	return b
}

func init() {
	MsgErrCodeMap = make(map[string]ErrCodeObj)
	MsgErrCodeMap[Ok] = ErrCodeObj{200, 20000, Ok}
	MsgErrCodeMap[OpErr] = ErrCodeObj{400, 40001, OpErr}
	MsgErrCodeMap[NotFound] = ErrCodeObj{404, 40004, NotFound}
	MsgErrCodeMap[NotExistBucket] = ErrCodeObj{404, 40005, NotExistBucket}
	MsgErrCodeMap[NetErr] = ErrCodeObj{409, 40002, NetErr}
	MsgErrCodeMap[InvalidKey] = ErrCodeObj{403, 40003, InvalidKey}
	MsgErrCodeMap[InvalidRequest] = ErrCodeObj{403, 40005, InvalidRequest}
	MsgErrCodeMap[QsInvalid] = ErrCodeObj{403, 40006, QsInvalid}
	MsgErrCodeMap[InternalErr] = ErrCodeObj{500, 50000, InternalErr}
}
