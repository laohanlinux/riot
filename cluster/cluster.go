package cluster

import (
	//"io/ioutil"
	"time"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/fsm"
	rstore "github.com/laohanlinux/riot/store"
	"os"
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

	// create log store dir, may be use disk
	store := raft.NewInmemStore()
	// rCluster.Dir = dir
	// for log and config storage
	rCluster.Stores = store

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
	r, err := raft.NewRaft(conf, rCluster.FSM, store, store, snap, peerStorage, tran)
	if err != nil {
		logger.Fatal(err)
	}
	rCluster.R = r

	//go rCluster.LeaderChange()
	return rCluster
}

func (c *Cluster) Join() {

}
func (c *Cluster) Status() string {
	return c.R.State().String()
}

func (c *Cluster) LeaderChange() {
	for {
		logger.Info("leader is: ", c.R.Leader())
		<-c.R.LeaderCh()
		logger.Info("leader change to ", c.R.Leader())
	}
}

func (c *Cluster) Leader() string {
	return c.R.Leader()
}

func (c *Cluster) Get(key string) ([]byte, error) {
	return c.FSM.Get(key)
}
