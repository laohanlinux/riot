package controller

import (
	"fmt"
	"strconv"

	"github.com/laohanlinux/riot/proxy/clientrpc"
	"github.com/laohanlinux/riot/proxy/http/errcode"
	"github.com/laohanlinux/riot/proxy/http/middleware"

	log "github.com/laohanlinux/utils/gokitlog"
	"gopkg.in/macaron.v1"
)

// TODO
// consistent
func GetValue(ctx *macaron.Context) {
	var (
		key     = ctx.Params("key")
		bucket  = ctx.Params("bucket")
		qsValue = ctx.Req.URL.Query().Get("qs")
		qs      bool
		res, _  = ctx.Data[middleware.ResKey].(map[string]interface{})
		value   []byte
		has     bool
		err     error
	)
	if qsValue != "" {
		if qs, err = strconv.ParseBool(qsValue); err != nil {
			log.Error("err", err)
			res["ret"] = errcode.ErrCodeInvalidRequest
			return
		}
	}

	if value, has, err = clientrpc.KV(bucket, key, qs); err != nil {
		log.Error("err", err)
		return
	}
	if !has {
		res["ret"] = errcode.ErrCodeNotFound
		return
	}
	res["data"] = fmt.Sprintf("%s", value)
}

func SetValue(ctx *macaron.Context) {
	var (
		key        = ctx.Params("key")
		bucket     = ctx.Params("bucket")
		res, _     = ctx.Data[middleware.ResKey].(map[string]interface{})
		value, err = ctx.Req.Body().Bytes()
	)
	if err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInvalidRequest
		return
	}
	if err = clientrpc.SetKV(bucket, key, value); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
}

func DelValue(ctx *macaron.Context) {
	var (
		key    = ctx.Params("key")
		bucket = ctx.Params("bucket")
		res, _ = ctx.Data[middleware.ResKey].(map[string]interface{})
		err    error
	)
	if err = clientrpc.DelKey(bucket, key); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
}
