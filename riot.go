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


func main() {
	var cfgPath string
	var joinAddr string
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
			go join(cfg, joinAddr)
		}
		// rc.EnableSingleNode = true
		cluster.NewCluster(cfg, rc)
		updateShareMemory(cfg)
		m := mux.NewRouter()
		m.Handle("/riot", &handler.RiotHandler{})
		m.HandleFunc("/admin/{cmd}", handler.AdminHandlerFunc)
		if err := http.ListenAndServe(cfg.SC.Addr+":"+cfg.SC.Port, m); err != nil {
			fmt.Printf("%s\n", err)
		}

	}()

	gGroup.Wait()
}

func join(cfg *config.Configure, joinAddr string) {
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
		for {
			time.Sleep(time.Second * 3)

			// get leaderName
			r := cluster.SingleCluster().R
			if cfg.RaftC.AddrString() == r.Leader() {
				// update leader addr info
				opRequest := pb.OpRequest{
					Op:    command.CmdShare,
					Key:   "",
					Value: []byte(cfg.RpcC.AddrString()),
				}
				b, _ := json.Marshal(opRequest)
				err := r.Apply(b, 3)
				if err != nil && err.Error() != nil {
					logger.Warn(r.Leader(), err.Error())
					continue
				}
				time.Sleep(time.Second * 5)
			}
			cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port = share.ShCache.LRPC.Addr, share.ShCache.LRPC.Port
		}
	}()
}
