package cluster

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/fsm"
	rstore "github.com/laohanlinux/riot/store"
)

type Cluster struct {
	Dir         string
	R           *raft.Raft
	Stores      *raft.InmemStore
	FSM         *fsm.StorageFSM
	Snap        raft.SnapshotStore
	Tran        raft.Transport
	PeerStorage raft.PeerStore
}

var rCluster *Cluster

// SingleCluster .
func SingleCluster() *Cluster {
	return rCluster
}

// NewCluster ...
func NewCluster(cfg *config.Configure, conf *raft.Config) *Cluster {
	if rCluster != nil {
		return rCluster
	}
	rCluster = &Cluster{}

	store := raft.NewInmemStore()
	rCluster.Stores = store
	// init raft applog and raftlog
	applyStore := initRaftLog(cfg, conf)
	// create a key/value store
	if err := os.RemoveAll(cfg.RaftC.StoreBackendPath); err != nil {
		logger.Fatal(err)
	}
	edbs := rstore.NewLeveldbStorage(cfg.RaftC.StoreBackendPath)
	rCluster.FSM = fsm.NewStorageFSM(edbs)

	//create snap dir
	snap, err := raft.NewFileSnapshotStore(cfg.RaftC.SnapshotStorage, 3, nil)
	if err != nil {
		logger.Fatal(err)
	}
	rCluster.Snap = snap

	// create transport
	tran, err := raft.NewTCPTransport(cfg.RaftC.Addr+":"+cfg.RaftC.Port, nil, 3, 2*time.Second, nil)
	if err != nil {
		logger.Fatal(err)
	}
	rCluster.Tran = tran

	// NewJSONPeers create a new JSONPees store
	peerStorage := raft.NewJSONPeers(cfg.RaftC.PeerStorage, tran)

	ps, err := peerStorage.Peers()
	if cfg.RaftC.EnableSingleNode && len(ps) <= 1 {
		conf.EnableSingleNode = cfg.RaftC.EnableSingleNode
		conf.DisableBootstrapAfterElect = false
	}

	peerStorage.SetPeers(ps)

	rCluster.PeerStorage = peerStorage
	// Wait the transport
	r, err := raft.NewRaft(conf, rCluster.FSM, applyStore, applyStore, snap, peerStorage, tran)
	if err != nil {
		logger.Fatal(err)
	}
	rCluster.R = r

	return rCluster
}

func (c *Cluster) Join() {

}
func (c *Cluster) Status() string {
	return c.R.State().String()
}

func (c *Cluster) LeaderChange(cfg *config.Configure) {
	go func() {
		for {
			time.Sleep(time.Second * 3)
			leaderName := c.R.Leader()
			logger.Debug(leaderName)
			if leaderName == "" {
				continue
			}
			oldLeaderName := cfg.LeaderRpcC.Addr + ":" + cfg.LeaderRpcC.Port
			logger.Debug(oldLeaderName)
			if oldLeaderName == leaderName {
				continue
			}
			addr := strings.Split(leaderName, ":")
			// url := fmt.Sprintf("http://%s:%s/admin/lrpc", addr[0], cfg.SC.Port)
			// logger.Info("lrpc request:", url)
			// resp, err := http.Get(url)
			// if err != nil {
			// 	logger.Error(err)
			// 	continue
			// }
			// rmsg := msgpack.ResponseMsg{}
			// b, err := ioutil.ReadAll(resp.Body)
			// if err != nil {
			// 	logger.Error(err)
			// 	resp.Body.Close()
			// 	continue
			// }
			// err = json.Unmarshal(b, &rmsg)
			// if err != nil {
			// 	logger.Error(err)
			// 	resp.Body.Close()
			// 	continue
			// }
			// resp.Body.Close()
			// logger.Info(rmsg)
			// addr = strings.Split(fmt.Sprintf("%s", rmsg.Results), ":")
			cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port = addr[0], cfg.RpcC.Port
			logger.Info("leader change to ", addr, leaderName)
		}
	}()
}

func (c *Cluster) Leader() string {
	return c.R.Leader()
}

func (c *Cluster) Get(key string) ([]byte, error) {
	return c.FSM.Get(key)
}

func initRaftLog(cfg *config.Configure, conf *raft.Config) *raftboltdb.BoltStore {
	// init raft app log
	logFile := path.Join(cfg.RaftC.RaftLogPath, "raft.log")
	fp, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	conf.LogOutput = fp
	// init raft apply log
	raftDBPath := path.Join(cfg.RaftC.ApplyLogPath, "apply.log")
	bdb, err := raftboltdb.NewBoltStore(raftDBPath)
	if err != nil {
		panic(err)
	}
	return bdb
}
