package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/cmd"
	"github.com/laohanlinux/riot/rpc"

	"github.com/boltdb/bolt"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/mux"
)

type RiotBucketHandler struct{}

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
	msg := MsgErrCodeMap[errType]
	if msg.httpCode == 200 {
		w.Write(value)
		return
	}

	w.WriteHeader(msg.httpCode)
	w.Write(msg.toJSONBytes())
}

func getBucket(w http.ResponseWriter, r *http.Request) (string, []byte, error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	rcmd := rpc.RpcCmd{
		Op:     cmd.CmdGetBucket,
		Bucket: bucket,
		Key:    "",
		Value:  nil,
	}

	if len(bucket) == 0 {
		return InvalidBucket, nil, fmt.Errorf("the request bucket is empty")
	}

	qs := cmd.QsConsistent

	value, err := rcmd.DoGet(qs)

	if err != nil && err.Error() != bolt.ErrBucketNotFound.Error() {
		logger.Error(err, bolt.ErrBucketNotFound)
		return OpErr, value, err
	}
	if err != nil && err.Error() == bolt.ErrBucketNotFound.Error() {
		return NotExistBucket, nil, nil
	}

	return Ok, value, nil
}

func setBucket(w http.ResponseWriter, r *http.Request) (string, error) {
	// vars := mux.Vars(r)
	//	bucket := vars["bucket"]
	value, err := ioutil.ReadAll(r.Body)
	if err != nil || value == nil || len(value) == 0 {
		return InvalidRequest, err
	}

	rcmd := rpc.RpcCmd{
		Op:     cmd.CmdCreateBucket,
		Bucket: string(value),
		Key:    "",
		Value:  nil,
	}

	err = rcmd.DoSet()
	if err != nil {
		return InternalErr, err
	}
	return Ok, nil
}

func delBucket(w http.ResponseWriter, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	cmd := rpc.RpcCmd{
		Op:     cmd.CmdDelBucket,
		Bucket: bucket,
		Key:    "",
		Value:  nil,
	}

	err := cmd.DoSet()
	if err != nil && err != cluster.ErrNotFound {
		return OpErr, err
	}

	if err == cluster.ErrNotFound {
		return NotFound, nil
	}

	return Ok, nil
}
