package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/laohanlinux/riot/command"
	"github.com/laohanlinux/riot/fsm"

	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/mux"
)

// ...
const (
	Ok             = "ok"
	OpErr          = "operation errror"
	NetErr         = "net work timeout"
	InternalErr    = "riot server error"
	InvalidKey     = "invalid key"
	InvalidRequest = "invalid Request"
	QsInvalid      = "invalid query strategies"
	NotFound       = "not found"
)

var msgErrCodeMap map[string]errCodeObj

type errCodeObj struct {
	httpCode   int
	StatusCode int    `json:"errCode"`
	Info       string `json:"msg"`
}

func (err *errCodeObj) toJsonBytes() []byte {
	b, _ := json.Marshal(err)
	return b
}

func init() {
	msgErrCodeMap = make(map[string]errCodeObj)
	msgErrCodeMap[Ok] = errCodeObj{200, 20000, Ok}
	msgErrCodeMap[OpErr] = errCodeObj{400, 40001, OpErr}
	msgErrCodeMap[NotFound] = errCodeObj{404, 40004, NotFound}
	msgErrCodeMap[NetErr] = errCodeObj{409, 40002, NetErr}
	msgErrCodeMap[InvalidKey] = errCodeObj{403, 40003, InvalidKey}
	msgErrCodeMap[InvalidRequest] = errCodeObj{403, 40005, InvalidRequest}
	msgErrCodeMap[QsInvalid] = errCodeObj{403, 40006, QsInvalid}
	msgErrCodeMap[InternalErr] = errCodeObj{500, 50000, InternalErr}
}

// RiotHandler ...
type RiotHandler struct{}

// ServeHTTP .
func (rh *RiotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var value []byte
	var err error
	var errType string

	switch r.Method {
	case "GET":
		errType, value, err = getValue(w, r)
		if err != nil {
			logger.Error(err)
		}
	case "DELETE":
		errType, err = delValue(w, r)
		if err != nil {
			logger.Error(err)
		}
	case "POST":
		errType, err = setValue(w, r)
		if err != nil {
			logger.Error(errType, err)
		}
	default:
		errType = InvalidRequest
	}
	msg := msgErrCodeMap[errType]
	if msg.httpCode == 200 {
		w.Write(value)
		return
	}
	w.WriteHeader(msg.httpCode)
	w.Write(msg.toJsonBytes())
}

func getValue(w http.ResponseWriter, r *http.Request) (string, []byte, error) {
	vars := mux.Vars(r)
	key := vars["key"]

	cmd := command.Command{
		Op:  command.CmdGet,
		Key: key,
	}
	if len(cmd.Key) == 0 {
		return InvalidKey, nil, fmt.Errorf("The Key is Empty")
	}

	qs := command.QsRandon
	var err error
	//Query strategires
	qsValue := r.URL.Query().Get("qs")
	if qsValue == "" {
		qs, err = strconv.Atoi(qsValue)
		if err != nil {
			return QsInvalid, nil, err
		}
	}
	var value []byte
	value, err = cmd.DoGet(qs)
	if err != nil && err != fsm.ErrNotFound {
		return OpErr, value, err
	}
	if err == fsm.ErrNotFound {
		return NotFound, nil, nil
	}
	return Ok, value, nil
}

func setValue(w http.ResponseWriter, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := ioutil.ReadAll(r.Body)
	if err != nil || value == nil || len(value) == 0 {
		return InvalidRequest, err
	}
	cmd := command.Command{
		Op:    command.CmdSet,
		Key:   key,
		Value: value,
	}
	err = cmd.DoSet()
	if err != nil {
		return InternalErr, err
	}
	return Ok, nil
}

func delValue(w http.ResponseWriter, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	key := vars["key"]
	cmd := command.Command{
		Op:  command.CmdDel,
		Key: key,
	}

	err := cmd.DoDel()
	if err != nil && err != fsm.ErrNotFound {
		return OpErr, err
	}

	if err == fsm.ErrNotFound {
		return NotFound, nil
	}
	return Ok, nil
}
