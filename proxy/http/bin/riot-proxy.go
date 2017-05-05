package main

import (
	"context"
	"flag"
	"net/http"
	"syscall"
	"time"

	"github.com/laohanlinux/riot/platform"
	"github.com/laohanlinux/riot/proxy/clientrpc"
	"github.com/laohanlinux/riot/proxy/http/config"
	"github.com/laohanlinux/riot/proxy/http/router"
	log "github.com/laohanlinux/utils/gokitlog"
)

func main() {
	var (
		configFile string
		s          http.Server
		c          = make(chan error)
		err        error
	)

	flag.StringVar(&configFile, "c", "cfg.toml", "")
	flag.Parse()
	config.InitConfig(configFile)

	// init rpc
	if err = clientrpc.InitRPC(config.Conf.Riot.RpcAddr, config.Conf.Riot.PoolSize); err != nil {
		panic(err)
	}

	// init http
	s = http.Server{
		Addr:              config.Conf.Srv.RPCAddr,
		ReadHeaderTimeout: time.Second * 3,
		WriteTimeout:      time.Second * 3,
		MaxHeaderBytes:    1 << 20,
	}
	if config.Conf.Srv.Token != "" {
		s.Handler = router.NewRouter(config.Conf.Srv.Token)
	} else {
		s.Handler = router.NewRouter()
	}
	defer s.Shutdown(context.Background())
	go func() {
		log.Debug("riot-proxy http server", config.Conf.Srv.RPCAddr)
		c <- s.ListenAndServe()
	}()

	platform.RegistSignal(syscall.SIGINT)
	log.Debug("server exit", <-c)
}
