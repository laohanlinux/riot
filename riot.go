package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/platform"
	"github.com/laohanlinux/riot/rpc"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
)

func main() {
	var cfgPath, joinAddr string
	flag.StringVar(&cfgPath, "c", "", "configure path")
	flag.StringVar(&joinAddr, "join", "", "host:port of leader to join")
	flag.Parse()
	var action = flag.Arg(0)

	runtime.GOMAXPROCS(runtime.NumCPU())

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	if cfgPath == "" {
		fmt.Println("No config path")
		return
	}

	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg, err := config.NewConfig(string(data))
	if err != nil {
		fmt.Println(err)
	}

	// Init log confiure
	logger.SetConsole(true)
	err = logger.SetRollingDaily(cfg.LogC.LogDir, cfg.LogC.LogName)
	if err != nil {
		panic(err)
	}

	var gGroup sync.WaitGroup
	var rpcService rpc.RiotRPCService
	// Init rpc server
	gGroup.Add(1)
	go func() {
		defer gGroup.Done()
		var err error
		rpcService, err = rpc.NewRpcServer(cfg.RpcC.Addr + ":" + cfg.RpcC.Port)
		if err != nil {
			panic(err)
		}
		logger.Info("Start rpc server successfully")
	}()
	gGroup.Add(1)
	go func() {
		defer gGroup.Done()
		// Init raft server
		rc := raft.DefaultConfig()
		// set snapshot
		if action == "dev" {
			rc.SnapshotThreshold = 10
			rc.TrailingLogs = 10
		}
		cluster.NewCluster(cfg, rc)
		// register monitor
		if action == "dev" {
			go func() {
				http.ListenAndServe(cfg.SMC.Addr+":"+cfg.SMC.Port, nil)
			}()
		}
	}()

	// regist the signal
	platform.RegistSignal(syscall.SIGINT)

	gGroup.Wait()
}
