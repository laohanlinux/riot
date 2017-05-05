package cluster

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/store"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	log "github.com/laohanlinux/utils/gokitlog"
)

type Cluster struct {
	Dir         string
	R           *raft.Raft
	Stores      *raft.InmemStore
	FSM         *StorageFSM
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

	memStore := raft.NewInmemStore()
	rCluster.Stores = memStore
	// init raft applog and raftlog
	applyStore := initRaftLog(cfg, conf)
	// create a key/value store
	if err := os.RemoveAll(cfg.RaftC.StoreBackendPath); err != nil {
		log.Crit("err", err)
	}
	edbs := store.NewRiotStoreFactory(cfg.RaftC.StoreBackend, cfg.RaftC.StoreBackendPath)
	rCluster.FSM = NewStorageFSM(edbs)

	//create snap dir
	snap, err := raft.NewFileSnapshotStore(cfg.RaftC.SnapshotStorage, 3, nil)
	if err != nil {
		log.Crit("err", err)
	}
	rCluster.Snap = snap

	// create transport
	tranAddr := fmt.Sprintf("%s:%s", cfg.RaftC.Addr, cfg.RaftC.Port)
	tran, err := raft.NewTCPTransport(tranAddr, nil, 3, 2*time.Second, nil)
	if err != nil {
		log.Crit("err", err)
	}
	rCluster.Tran = tran

	// NewJSONPeers create a new JSONPees store
	peerStorage := raft.NewJSONPeers(cfg.RaftC.PeerStorage, tran)

	ps, err := peerStorage.Peers()
	if cfg.RaftC.EnableSingleNode && len(ps) <= 1 {
		conf.EnableSingleNode = cfg.RaftC.EnableSingleNode
		conf.DisableBootstrapAfterElect = false
	}
	// get the peers
	ps = cfg.RaftC.Peers

	peerStorage.SetPeers(ps)

	rCluster.PeerStorage = peerStorage
	// Wait the transport
	log.Debug("waitting for electing leader...")
	r, err := raft.NewRaft(conf, rCluster.FSM, applyStore, applyStore, snap, peerStorage, tran)
	log.Debug("has elected the leader.")
	if err != nil {
		log.Crit(err)
	}
	rCluster.R = r
	peers, _ := rCluster.PeerStorage.Peers()
	log.Info("addr", rCluster.Tran.LocalAddr(), "status", r.State().String(), "peers", fmt.Sprintf("%+v", peers))

	return rCluster
}

// Status of the node running info
func (c *Cluster) Status() string {
	return c.R.State().String()
}

func (c *Cluster) LocalString() string {
	return c.R.String()
}

func (c *Cluster) Leader() string {
	return c.R.Leader()
}

func (c *Cluster) Get(bucket, key []byte) ([]byte, error) {
	return c.FSM.Get(bucket, key)
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
