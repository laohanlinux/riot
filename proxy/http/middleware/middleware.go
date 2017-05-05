package middleware

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/laohanlinux/riot/proxy/http/errcode"
	log "github.com/laohanlinux/utils/gokitlog"
	macaron "gopkg.in/macaron.v1"
)

var ResKey = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("=======%v========", time.Now().Unix()))))

func ElaspRequest(ctx *macaron.Context) {
	var timeStart = time.Now()
	ctx.Next()
	log.Infof("method: %v, statusCode:%v req: %s, ip: %s, time: %fs", ctx.Req.Method, ctx.Resp.Status(), ctx.Req.URL.String(), ctx.Req.RemoteAddr, time.Now().Sub(timeStart).Seconds())
}

func ContextMiddleware(ctx *macaron.Context) {
	var (
		res = map[string]interface{}{"ret": errcode.Ok}
	)
	ctx.Data[ResKey] = res
	ctx.Next()
}

func AuthorMiddleware(token string, ctx *macaron.Context) {
	var (
		res = ctx.Data[ResKey].(map[string]interface{})
	)
	if token != "" && token != strings.ToLower(ctx.Req.Header.Get("x-token")) {
		res["ret"] = errcode.ErrCodeForbidden
		OutputMiddleware(ctx)
		return
	}
	ctx.Next()
}

func OutputMiddleware(ctx *macaron.Context) {
	var (
		res, ok = ctx.Data[ResKey].(map[string]interface{})
		errCode int
		err     error
	)
	ctx.Next()
	if !ok {
		return
	}
	if ctx.Resp.Written() {
		return
	}
	ctx.Resp.Header().Add("Content-Type", "application/json")
	errCode = res["ret"].(int)
	switch errCode {
	case errcode.Ok:
		ctx.Resp.WriteHeader(http.StatusOK)
	case errcode.ErrCodeInternal:
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
	case errcode.ErrCodeInvalidRequest, errcode.ErrCodeForbidden:
		ctx.Resp.WriteHeader(http.StatusForbidden)
	case errcode.ErrCodeNotFound:
		ctx.Resp.WriteHeader(http.StatusNotFound)
	default:
		ctx.Resp.WriteHeader(http.StatusOK)
	}

	if err = json.NewEncoder(ctx.Resp).Encode(res); err != nil {
		log.Error("err", err)
	}
}
