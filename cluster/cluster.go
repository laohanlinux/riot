package cluster

import (
	"io/ioutil"
	"net"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/fsm"
)

type Cluster interface {
	//NewCluster(*raft.Config)
	Join(string) error
	Leave() error
	Close() error
	RemoveNode(string) error
	IsLeader() bool
	LeaderCh() <-chan bool
	Nodes() ([]string, error)
}

type riotCluster struct {
	peerAddres []string
	n          node
}

func NewCluster(localAddr string, peerAddres []string, conf *raft.Config) *Cluster {
	rcluster := riotCluster{
		peerAddres: make([]string, 0),
	}

	if a, err := net.ResolveTCPAddr("tcp", localAddr); err != nil {
		panic(err)
	} else {
		rcluster.n.addr = localAdd
	}

	perrs := make([]string, 0)
	for _, addr := range peerAddres {
		if p, err := net.ResolveTCPAddr("tcp", addr); err != nil {
			panic(err)
		} else {
			perrs = raft.AddUniquePeer(perrs, p.String())
		}
	}

	rcluster.peerAddres = perrs
	// Setup the restores and transports
	for i := 0; i < n; i++ {
		dir, err := ioutil.TempDir("", "raft")
		if err != nil {
			logger.Fatal(err)
		}

		// create log store dir
		store := raft.NewInmemStore()
		rcluster.n.dir = dir
		rcluster.n.stores = store
		rcluster.n.fsm = &fsm.StorageFSM{}

		//create snap dir
		_, snap := fileSnap()
		rcluster.n.snap = snap

		// create transport
		addr, tran := raft.NewTCPTransport(rcluster)
		rcluster.n.tran = tran
		rcluster.localAdd = addr
	}

	// Wait the transport

	// Create peerStore from log data
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
