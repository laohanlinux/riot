package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/command"
)

// ...
const (
	OpErr       = 0
	NetErr      = 1
	InternalErr = 2
	InvalidErr  = 3
	InvalidKey  = 4
)

type msgErrCode struct {
	ErrCode int `json:"errCode"`
}

func (m *msgErrCode) setJson(errCode ...int) []byte {
	if len(errCode) > 0 {
		m.ErrCode = errCode[0]
	}
	b, _ := json.Marshal(m)
	return b
}

// RiotHandler ...
type RiotHandler struct {
}

// ServeHTTP .
func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mErrCode := &msgErrCode{
		ErrCode: 0,
	}

	defer func() {
		if err := recover(); err != nil {
			mErrCode.ErrCode = InternalErr
			w.WriteHeader(500)
		}
		w.Write(mErrCode.setJson())
	}()

	log.Info("RiotHandler receive a request:", r.URL.Path)

	value, err := ioutil.ReadAll(r.Body)

	if err != nil || value == nil {
		mErrCode.ErrCode = InvalidErr
		return
	}

	switch r.Method {
	case "GET":

	case "DELETE":
	case "POST":
	default:
	}

}

func getValue(w http.ResponseWriter, r *http.Request, mErrCode *msgErrCode) {
	cmd := command.Command{
		Op:  command.CmdGet,
		Key: r.URL.RequestURI(),
	}
	if len(cmd.Key) == 0 {
		mErrCode.ErrCode = InvalidKey
		return
	}

	cmd.DoGet()
}

func setValue(w http.ResponseWriter, r *http.Request) {

}

func delValue(w http.ResponseWriter, r *http.Request) {

}
