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
	"github.com/laohanlinux/riot/command"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/handler"
	"github.com/laohanlinux/riot/handler/msgpack"
	"github.com/laohanlinux/riot/rpc"
	"github.com/laohanlinux/riot/rpc/pb"
	"github.com/laohanlinux/riot/share"

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
		c := cluster.NewCluster(cfg, rc)
		updateShareMemory(cfg)
		// waitting for leader
		c.LeaderChange(cfg)
		m := mux.NewRouter()
		m.Handle("/riot", &handler.RiotHandler{})
		m.HandleFunc("/admin/{cmd}", handler.AdminHandlerFunc)
		if err := http.ListenAndServe(cfg.SC.Addr+":"+cfg.SC.Port, m); err != nil {
			fmt.Printf("%s\n", err)
		}

	}()

	gGroup.Wait()
}

func join(cfg *config.Configure) {
	var results msgpack.ResponseMsg
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

func updateShareMemory(cfg *config.Configure) {

	go func() {
		var sc share.ShareCache
		sc.LRPC = &share.LeaderRpcAddr{
			Addr: "",
			Port: "",
		}

		for {
			time.Sleep(time.Second * 3)
			oddr, oport := share.ShCache.LRPC.Addr, share.ShCache.LRPC.Port
			oddr1, oprt1 := cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port
			if (oddr+oport) == (oddr1+oprt1) && oddr != "" {
				continue
			}
			// old not equall new or old is empty
			logger.Info(cfg.LeaderRpcC, string(share.ShCache.ToBytes()))
			// update data by leader node
			cfg.LeaderRpcC.Addr = oddr
			cfg.LeaderRpcC.Port = oport
			leaderAddr := cluster.SingleCluster().R.Leader()
			cfg.LeaderRpcC.Addr = oddr
			cfg.LeaderRpcC.Port = oport

			if leaderAddr != (cfg.RaftC.Addr + ":" + cfg.RaftC.Port) {
				logger.Info("不是leader, leader:", leaderAddr, cfg.RaftC.Addr+":"+cfg.RaftC.Port)
				continue
			}
			// encode ShareCache
			sc.LRPC.Addr = cfg.RpcC.Addr
			sc.LRPC.Port = cfg.RpcC.Port
			opRequest := pb.OpRequest{
				Op:    command.CmdShare,
				Key:   "",
				Value: sc.ToBytes(),
			}
			b, _ := json.Marshal(opRequest)
			err := cluster.SingleCluster().R.Apply(b, 3)
			if err != nil {
				logger.Error(err)
			}
		}
	}()
}
