package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/laohanlinux/riot/command"
	"github.com/laohanlinux/riot/fsm"

	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/mux"
)

type RiotBucketHandler struct {
}

func (rbh *RiotBucketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var value []byte
	var err error
	var errType string
	w.Header().Add("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		errType, value, err = getBucket(w, r)
		if err != nil {
			logger.Error(err)
		}
	case "DELETE":
		errType, err = delBucket(w, r)
		if err != nil {
			logger.Error(err)
		}
	case "POST":
		errType, err = setBucket(w, r)
		if err != nil {
			logger.Error(err)
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

func getBucket(w http.ResponseWriter, r *http.Request) (string, []byte, error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	cmd := command.Command{
		Op:     command.CmdGetBucket,
		Bucket: bucket,
		Key:    nil,
		Value:  nil,
	}

	if len(bucket) == 0 {
		return InvalidBucket, nil, fmt.Errorf("the request bucket is empty")
	}

	qs := command.QsConsistent

	var value []byte
	value, err = cmd.DoGet(qs)
	if err != nil && err != fsm.ErrNotFound {
		return OpErr, value, err
	}
	if err == fsm.ErrNotFound {
		return NotExistBucket, nil, nil
	}

	return Ok, value, nil
}

func setBucket(w http.ResponseWriter, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	value, err := ioutil.ReadAll(r.Body)
	if err != nil || value == nill || len(value) == 0 {
		return InvalidRequest, err
	}

	cmd := command.CmdDelBucket{
		Op:     command.CmdSetBucket,
		bucket: bucket,
		Key:    nil,
		Value:  nil,
	}

	err = cmd.DoSet()
	if err != nil {
		return InternalErr, err
	}
	return Ok, nil
}

func delBucket(w http.ResponseWriter, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	cmd := command.Command{
		Op:     command.CmdDelBucket,
		Bucket: bucket,
		Key:    nil,
		Value:  nil,
	}

	err := cmd.DoSet()
	if err != nil && err != fsm.ErrNotFound {
		return OpErr, err
	}

	if err == fsm.ErrNotFound {
		return NotFound, nil
	}

	return Ok, nil
}
