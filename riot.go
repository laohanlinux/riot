package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"

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

	// Init rpc server
	go func() {
		_, err := rpc.NewRpcServer(cfg.RpcC.Addr + ":" + cfg.RpcC.Port)
		if err != nil {
			panic(err)
		}
		fmt.Println("Start rpc server successfully")
	}()

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
		b, _ := json.Marshal(map[string]string{"ip": cfg.RaftC.Addr, "port": cfg.RaftC.Port})
		resp, err := http.Post("http://"+joinAddr+"/admin/join", "application-type/json", bytes.NewReader(b))
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Printf("%v\n", resp)
	}
	// rc.EnableSingleNode = true
	cluster.NewCluster(cfg, rc)

	m := mux.NewRouter()
	m.Handle("/riot", &handler.RiotHandler{})
	m.HandleFunc("/admin/{cmd}", handler.AdminHandlerFunc)
	if err := http.ListenAndServe(cfg.SC.Addr+":"+cfg.SC.Port, m); err != nil {
		fmt.Printf("%s\n", err)
	}
}
