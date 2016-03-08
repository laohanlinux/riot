package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/laohanlinux/riot/command"
	"github.com/laohanlinux/riot/fsm"
)

// ...
const (
	Ok             = "ok"
	OpErr          = "operation errror"
	NetErr         = "net work timeout"
	InternalErr    = "riot server error"
	InvalidKey     = "invalid key"
	InvalidRequest = "Invalid Request"
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
			fmt.Printf("%s\n", err)
		}
	case "DELETE":
		errType, err = delValue(w, r)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	case "POST":
		errType, err = setValue(w, r)
		if err != nil {
			fmt.Printf("%s\n", err)
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
	cmd := command.Command{
		Op:  command.CmdGet,
		Key: r.URL.RequestURI(),
	}
	if len(cmd.Key) == 0 {
		return InvalidKey, nil, fmt.Errorf("The Key is Empty")
	}

	value, err := cmd.DoGet()
	if err != nil && err != fsm.ErrNotFound {
		return OpErr, value, err
	}
	if err == fsm.ErrNotFound {
		return NotFound, nil, nil
	}
	return Ok, value, nil
}

func setValue(w http.ResponseWriter, r *http.Request) (string, error) {
	value, err := ioutil.ReadAll(r.Body)
	if err != nil || value == nil {
		return InvalidRequest, err
	}
	cmd := command.Command{
		Op:    command.CmdSet,
		Key:   r.URL.RequestURI(),
		Value: value,
	}
	err = cmd.DoSet()
	if err != nil {
		return InternalErr, err
	}
	return Ok, nil
}

func delValue(w http.ResponseWriter, r *http.Request) (string, error) {
	cmd := command.Command{
		Op:  command.CmdDel,
		Key: r.URL.RequestURI(),
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
