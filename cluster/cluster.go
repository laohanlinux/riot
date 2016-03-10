package cluster

import (
	"io/ioutil"
	"time"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/fsm"
)

type Cluster struct {
	n *Node
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
	rCluster = &Cluster{
		n: &Node{},
	}

	// Setup the restores and transports
	dir, err := ioutil.TempDir("", "raft")
	if err != nil {
		logger.Fatal(err)
	}

	// create log store dir, may be use disk
	store := raft.NewInmemStore()
	rCluster.n.dir = dir
	// for log and config storage
	rCluster.n.stores = store
	rCluster.n.fsm = fsm.NewStorageFSM()

	//create snap dir
	_, snap := fileSnap()
	rCluster.n.snap = snap

	// create transport
	tran, err := raft.NewTCPTransport(cfg.RaftC.Addr+":"+cfg.RaftC.Port, nil, 3, 2*time.Second, nil)
	if err != nil {
		logger.Fatal(err)
	}
	rCluster.n.tran = tran

	// NewJSONPeers create a new JSONPees store
	peerStorage := raft.NewJSONPeers(cfg.RaftC.PeerStorage, tran)
	ps, err := peerStorage.Peers()
	if err != nil {
		logger.Fatal(err)
	}
	for _, peer := range cfg.RaftC.Peers {
		ps = raft.AddUniquePeer(ps, peer)
	}
	if cfg.RaftC.EnableSingleNode && len(ps) > 0 {
		conf.EnableSingleNode = false
	}
	peerStorage.SetPeers(ps)
	// Wait the transport
	r, err := raft.NewRaft(conf, rCluster.n.fsm, store, store, snap, peerStorage, tran)
	if err != nil {
		logger.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	// future := r.VerifyLeader()
	// if err := future.Error(); err != nil {
	// 	logger.Error(err)
	// }

	// for _, peer := range ps {
	// 	if peer != tran.LocalAddr() {
	// 		future := r.AddPeer(peer)
	// 		if err := future.Error(); err != nil {
	// 			logger.Error(err)
	// 		}
	// 	}
	// }

	rCluster.n.r = r
	go func() {
		for {
			time.Sleep(time.Second * 3)
			logger.Info(r.State())
		}
	}()
	return rCluster
}

func (c *Cluster) Join() {

}

func (c *Cluster) Node() *Node {
	return c.n
}

func (c *Cluster) Leader() string {
	return c.n.Leader()
}
func fileSnap() (string, *raft.FileSnapshotStore) {
	dir, err := ioutil.TempDir("", "raft")
	if err != nil {
		panic(err)
	}

	logger.Info("snap save dir:", dir)
	snap, err := raft.NewFileSnapshotStore(dir, 3, nil)
	if err != nil {
		panic(err)
	}

	return dir, snap
}
