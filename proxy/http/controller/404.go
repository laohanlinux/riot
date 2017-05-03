package controller

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

func Contr404(ctx *macaron.Context) {
	ctx.WriteHeader(http.StatusNotFound)
	ctx.Resp.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<img src="https://raw.githubusercontent.com/laohanlinux/riot/master/doc/riot.jpg">
<body>
</body>
</html>`))
}
