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

var msgErrCodeMap map[string]errCodeObj

type errCodeObj struct {
	httpCode   int
	StatusCode int    `json:"errCode"`
	Info       string `json:"msg"`
}

func (err *errCodeObj) toJSONBytes() []byte {
	b, _ := json.Marshal(err)
	return b
}

func init() {
	msgErrCodeMap = make(map[string]errCodeObj)
	msgErrCodeMap[Ok] = errCodeObj{200, 20000, Ok}
	msgErrCodeMap[OpErr] = errCodeObj{400, 40001, OpErr}
	msgErrCodeMap[NotFound] = errCodeObj{404, 40004, NotFound}
	msgErrCodeMap[NotExistBucket] = errCodeObj{404, 40005, NotExistBucket}
	msgErrCodeMap[NetErr] = errCodeObj{409, 40002, NetErr}
	msgErrCodeMap[InvalidKey] = errCodeObj{403, 40003, InvalidKey}
	msgErrCodeMap[InvalidRequest] = errCodeObj{403, 40005, InvalidRequest}
	msgErrCodeMap[QsInvalid] = errCodeObj{403, 40006, QsInvalid}
	msgErrCodeMap[InternalErr] = errCodeObj{500, 50000, InternalErr}
}
