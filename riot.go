package main

import (
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"syscall"

	"github.com/laohanlinux/riot/api"
	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/platform"

	"github.com/hashicorp/raft"
	log "github.com/laohanlinux/utils/gokitlog"
	"github.com/laohanlinux/utils/netrpc"
)

func main() {
	var (
		cfgPath string
		err     error
		data    []byte
		cfg     *config.Configure
		rc      = raft.DefaultConfig()
		c       *cluster.Cluster
	)
	flag.StringVar(&cfgPath, "c", "", "configure path")
	flag.Parse()
	var action = flag.Arg(0)

	runtime.GOMAXPROCS(runtime.NumCPU())

	if cfgPath == "" {
		panic("No config path")
	}

	if data, err = ioutil.ReadFile(cfgPath); err != nil {
		panic(err)
	}
	if cfg, err = config.NewConfig(string(data)); err != nil {
		panic(err)
	}

	var gGroup sync.WaitGroup

	// Init raft server
	// set snapshot
	if action == "dev" {
		rc.SnapshotThreshold = 10
		rc.TrailingLogs = 10
	}
	c = cluster.NewCluster(cfg, rc)
	// register monitor
	if action == "dev" {
		go func() {
			http.ListenAndServe(cfg.SMC.Addr+":"+cfg.SMC.Port, nil)
		}()
	}

	// Init rpc server
	gGroup.Add(1)
	go func() {
		defer gGroup.Done()
		addr := cfg.RpcC.Addr + ":" + cfg.RpcC.Port
		service := api.NewAPIService(api.NewMiniAPI(c), api.NewAdmAPI(c))
		l, err := net.Listen("tcp", addr)
		if err != nil {
			panic(err)
		}
		defer l.Close()
		srv := netrpc.NewServer()
		srv.Register(service)
		srv.Register(&netrpc.HealthCheck{})
		log.Info("Start rpc server successfully")
		srv.Accept(l)
	}()

	// regist the signal
	platform.RegistSignal(syscall.SIGINT)

	gGroup.Wait()
}
