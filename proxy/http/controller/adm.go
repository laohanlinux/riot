package controller

import (
	"encoding/json"
	"fmt"

	"github.com/laohanlinux/riot/proxy/clientrpc"
	"github.com/laohanlinux/riot/proxy/http/errcode"
	"github.com/laohanlinux/riot/proxy/http/middleware"

	"github.com/hashicorp/raft"
	log "github.com/laohanlinux/utils/gokitlog"
	macaron "gopkg.in/macaron.v1"
)

// TODO
// add join, remove

func Leader(ctx *macaron.Context) {
	var (
		res, _ = ctx.Data[middleware.ResKey].(map[string]interface{})
		leader string
		err    error
	)
	if leader, err = clientrpc.Leader(); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}

	res["data"] = leader
}

func States(ctx *macaron.Context) {
	var (
		res, _ = ctx.Data[middleware.ResKey].(map[string]interface{})
		states []string
		err    error
	)
	if states, err = clientrpc.States(); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
	res["data"] = states
}

func Peers(ctx *macaron.Context) {
	var (
		res, _ = ctx.Data[middleware.ResKey].(map[string]interface{})
		peers  []string
		err    error
	)
	if peers, err = clientrpc.Peers(); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
	res["data"] = peers
}

func SnapshotInfo(ctx *macaron.Context) {
	var (
		res, _  = ctx.Data[middleware.ResKey].(map[string]interface{})
		snapLen int
		err     error
	)
	if snapLen, err = clientrpc.Snapshot(); err != nil {
		if err.Error() != raft.ErrNothingNewToSnapshot.Error() {
			log.Error("err", err)
			res["ret"] = errcode.ErrCodeInternal
			return
		}
		log.Debug("snapshot", err)
	}
	res["data"] = snapLen
	return
}

func RemovePeer(ctx *macaron.Context) {
	var (
		res, _     = ctx.Data[middleware.ResKey].(map[string]interface{})
		remoteAddr map[string]interface{}
		addr       string
		err        error
	)
	if err = json.NewDecoder(ctx.Req.Body().ReadCloser()).Decode(&remoteAddr); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInvalidRequest
		return
	}
	addr = fmt.Sprintf("%v:%v", remoteAddr["ip"], remoteAddr["port"])
	if err = clientrpc.RemovePeer(addr); err != nil {
		log.Error("err", err)
		res["ret"] = errcode.ErrCodeInternal
		return
	}
}

func RouterTest(ctx *macaron.Context) {
	var (
		res, _ = ctx.Data[middleware.ResKey].(map[string]interface{})
		ll     = ctx.Params("ll")
	)
	log.Warn("ctx..................", ll)
	res["data"] = ll
	return
}
