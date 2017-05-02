package http

import (
	"crypto/md5"
	"fmt"
	"time"

	log "github.com/laohanlinux/utils/gokitlog"
	macaron "gopkg.in/macaron.v1"
)

var ResKey = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("=======%v========", time.Now().Unix()))))

func ElaspRequest(ctx *macaron.Context) {
	var timeStart = time.Now()
	defer func() {
		log.Debugf("request elasp time:%v", time.Now().Sub(timeStart).Seconds())
	}()
	ctx.Next()
}
