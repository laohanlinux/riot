package controller

import (
	"github.com/laohanlinux/riot/proxy/http/errcode"
	"github.com/laohanlinux/riot/proxy/http/middleware"

	"github.com/boltdb/bolt"
	"github.com/laohanlinux/riot/proxy/clientrpc"
	log "github.com/laohanlinux/utils/gokitlog"
	macaron "gopkg.in/macaron.v1"
)

func CreateBucket(ctx *macaron.Context) {
	var (
		bucketName, err = ctx.Req.Body().String()
		res, _          = ctx.Data[middleware.ResKey].(map[string]interface{})
	)
	if err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}

	if err = clientrpc.CreateBucket(bucketName); err != nil {
		if err.Error() == bolt.ErrBucketExists.Error() {
			return
		}
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
}

func DelBucket(ctx *macaron.Context) {
	var (
		res, _     = ctx.Data[middleware.ResKey].(map[string]interface{})
		bucketName = ctx.Params("bucket")
		err        error
	)
	if bucketName == "" {
		res["ret"] = errcode.ErrCodeInvalidRequest
		return
	}
	if err = clientrpc.DelBucket(bucketName); err != nil {
		if err.Error() == bolt.ErrBucketNotFound.Error() {
			return
		}
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
}

// TODO
// read state
func BucketInfo(ctx *macaron.Context) {
	var (
		res, _     = ctx.Data[middleware.ResKey].(map[string]interface{})
		bucketName = ctx.Params("bucket")
		err        error
		info       interface{}
		has        bool
	)
	if bucketName == "" {
		res["ret"] = errcode.ErrCodeInvalidRequest
		return
	}
	if info, has, err = clientrpc.BucketInfo(bucketName); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
	if !has {
		res["ret"] = errcode.ErrCodeNotFound
		return
	}
	res["data"] = info
	return
}
