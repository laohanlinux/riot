package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
func (rh *RiotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mErrCode := &msgErrCode{
		ErrCode: 0,
	}

	defer func() {
		fmt.Println("path:", r.URL.RawPath)
		fmt.Println("key:", r.URL.Path)
		if err := recover(); err != nil {
			mErrCode.ErrCode = InternalErr
			w.WriteHeader(500)
		}
		w.Write(mErrCode.setJson())
	}()

	switch r.Method {
	case "GET":
		errCode, value, err := getValue(w, r)
		if errCode > 0 {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("value is :%s\n", value)
		}
		mErrCode.ErrCode = errCode
	case "DELETE":
		errCode, err := delValue(w, r)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		mErrCode.ErrCode = errCode
	case "POST":
		errCode, err := setValue(w, r)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		mErrCode.ErrCode = errCode
	default:
	}

}

func getValue(w http.ResponseWriter, r *http.Request) (int, []byte, error) {
	cmd := command.Command{
		Op:  command.CmdGet,
		Key: r.URL.RequestURI(),
	}
	if len(cmd.Key) == 0 {
		return InvalidKey, nil, fmt.Errorf("The Key is Empty")
	}

	value, err := cmd.DoGet()
	if err != nil {
		return OpErr, value, err
	}
	return 0, value, nil
}

func setValue(w http.ResponseWriter, r *http.Request) (int, error) {
	value, err := ioutil.ReadAll(r.Body)
	if err != nil || value == nil {
		return InvalidErr, err
	}
	cmd := command.Command{
		Op:    command.CmdSet,
		Key:   r.URL.RequestURI(),
		Value: value,
	}
	err = cmd.DoSet()
	if err != nil {
		return OpErr, err
	}

	return 0, nil
}

func delValue(w http.ResponseWriter, r *http.Request) (int, error) {
	cmd := command.Command{
		Op:  command.CmdDel,
		Key: r.URL.RequestURI(),
	}

	err := cmd.DoDel()
	if err != nil {
		return OpErr, err
	}
	return 0, nil
}
