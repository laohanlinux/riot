package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/handler"
	"github.com/laohanlinux/riot/rpc"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/mux"
)

var joinAddr string

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "", "configure path")
	flag.StringVar(&joinAddr, "join", "", "host:port of leader to join")
	flag.Parse()
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
	cfg.Info()

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
		fmt.Println("Start rpc server successfully")
	}()

	gGroup.Add(1)
	go func() {
		defer gGroup.Done()
		// Ini log configure
		logger.SetConsole(true)
		err = logger.SetRollingDaily(cfg.LogC.LogDir, cfg.LogC.LogName)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Init raft server
		rc := raft.DefaultConfig()


		if joinAddr != "" {
			go join(cfg)
		}
		// rc.EnableSingleNode = true
		cluster.NewCluster(cfg, rc)

		m := mux.NewRouter()
		m.Handle("/riot", &handler.RiotHandler{})
		m.HandleFunc("/admin/{cmd}", handler.AdminHandlerFunc)
		if err := http.ListenAndServe(cfg.SC.Addr+":"+cfg.SC.Port, m); err != nil {
			fmt.Printf("%s\n", err)
		}

	}()

	gGroup.Wait()
	fmt.Println("raft has exited")
}

func join(cfg *config.Configure) {
	var results handler.ResponseMsg
	for {
		time.Sleep(time.Second)
		b, _ := json.Marshal(map[string]string{"ip": cfg.RaftC.Addr, "port": cfg.RaftC.Port})
		resp, err := http.Post("http://"+joinAddr+"/admin/join", "application-type/json", bytes.NewReader(b))
		if err != nil {
			logger.Error(err)
			continue
		}
		rpl, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
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
