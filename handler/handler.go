package handler

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	log "github.com/laohanlinux/go-logger/logger"
)

const (
	Op_Error       = 0
	Net_Error      = 1
	Internal_Error = 2
	Invalid_Error = 3
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

type RiotHandler struct {
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mErrCode := &msgErrCode{
		ErrCode: 0,
	}

	defer func() {
		if err := recover(); err != nil {
			mErrCode.ErrCode = Internal_Error
			w.WriteHeader(500)
		}
		w.Write(mErrCode.setJson())
	}()

	log.Info("RiotHandler receive a request:", r.URL.Path)

	key := r.URL.Path
	
	if len(key) <= 0 {

	}
	
	value, err :=ioutil.ReadAll(r.Body)
	
	if err != nil || value == nil {
		mErrCode.ErrCode = Invalid_Error
		return
	}
	
	log.Info("RiotHandler receive content size:", len(value))
	
}
