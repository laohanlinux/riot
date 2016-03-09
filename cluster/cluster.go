package cluster

import (
	"io/ioutil"
	"net"
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

	_, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		panic(err)
	}

	rCluster.n.addr = localAddr
	logger.Info(peerAddres)
	var peers []string
	for _, addr := range peerAddres {
		if p, err := net.ResolveTCPAddr("tcp", addr); err != nil {
			panic(err)
		} else {
			peers = raft.AddUniquePeer(peers, p.String())
		}
	}

	rCluster.peerAddres = peers
	logger.Debug(peers)
	rCluster.n.r.SetPeers(rCluster.peerAddres)

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
	tran, err := raft.NewTCPTransport(localAddr, nil, 3, 5*time.Second, nil)
	if err != nil {
		logger.Fatal(err)
	}
	rCluster.n.tran = tran

	// create peer storage
	peerStorage := raft.NewJSONPeers(cfg.RaftC.PeerStorage, tran)

	// Wait the transport
	r, err := raft.NewRaft(conf, rCluster.n.fsm, store, store, snap, peerStore, tran)
	if err != nil {
		logger.Fatal(err)
	}

	rCluster.n.r = r
	return rCluster
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
