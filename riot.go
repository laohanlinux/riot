package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/handler"
	"github.com/laohanlinux/riot/handler/msgpack"
	"github.com/laohanlinux/riot/platform"
	"github.com/laohanlinux/riot/rpc"
	"github.com/laohanlinux/riot/share"
	"github.com/laohanlinux/riot/store"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/mux"
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
		// set snapshot
		if joinAddr != "" {
			go join(cfg, joinAddr)
		}
		cluster.NewCluster(cfg, rc)
		// waitting the cluster init
		go share.UpdateShareMemory(cfg, cluster.SingleCluster().R)

		// set app http router
		m := mux.NewRouter()
		switch cfg.RaftC.StoreBackend {
		case store.LevelDBStoreBackend:
			m.Handle("/riot/key/{key}", &handler.RiotHandler{})
		case store.BoltDBStoreBackend:
			m.Handle("/riot/bucket", &handler.RiotBucketHandler{})
			m.Handle("/riot/bucket/{bucket}", &handler.RiotBucketHandler{})
			m.Handle("/riot/bucket/{bucket}/key/{key}", &handler.RiotHandler{})
		default:
			os.Exit(-1)
		}
		m.HandleFunc("/riot/admin/{cmd}", handler.AdminHandlerFunc)
		// register monitor
		if action == "dev" {
			go func() {
				http.ListenAndServe(cfg.SMC.Addr+":"+cfg.SMC.Port, nil)
			}()
		}

		if err := http.ListenAndServe(cfg.SC.Addr+":"+cfg.SC.Port, m); err != nil {
			logger.Error(err)
		}

	}()

	// regist the signal
	platform.RegistSignal(syscall.SIGINT)

	gGroup.Wait()
}

func join(cfg *config.Configure, joinAddr string) {
	var results msgpack.ResponseMsg
	for {
		time.Sleep(time.Second)
		b, _ := json.Marshal(map[string]string{"ip": cfg.RaftC.Addr, "port": cfg.RaftC.Port})
		httpURL := fmt.Sprintf("http://%s/riot/admin/join", joinAddr)
		resp, err := http.Post(httpURL, "application-type/json", bytes.NewReader(b))
		if err != nil {
			logger.Error(err)
			continue
		}
		if resp.StatusCode != 200 {
			logger.Error(httpURL, " request status code:", resp.Status)
			continue
		}
		rpl, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Warn(err)
			continue
		}
		resp.Body.Close()
		if err = json.Unmarshal(rpl, &results); err != nil {
			logger.Error(err)
			continue
		}
		logger.Info(results)
		return
	}
}
