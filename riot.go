package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/handler"
	"github.com/laohanlinux/riot/pb"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "", "configure path")
	flag.Parse()

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
		_, err := pb.NewRpcServer(cfg.RpcC.Addr + ":" + cfg.RpcC.Port)
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
	rc.EnableSingleNode = true
	cluster.NewCluster(cfg.RaftC.Addr+":"+cfg.RaftC.Port, nil, rc)

	if err := http.ListenAndServe(cfg.SC.Addr+":"+cfg.SC.Port, &handler.RiotHandler{}); err != nil {
		fmt.Printf("%s\n", err)
	}
}
