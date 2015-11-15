package cluster

import (
	"io/ioutil"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/fsm"
)

type Cluster interface {
	//NewCluster(*raft.Config)

}

type riotCluster struct {
	nodeAddrs []string
	n         node
}

func NewCluster(conf *raft.Config) *Cluster {
	rcluster := riotCluster{}
	peers := make([]string, 0, n)

	// Setup the restores and transports

	for i := 0; i < n; i++ {
		dir, err := ioutil.TempDir("", "raft")
		if err != nil {
			logger.Fatal(err)
		}

		store := raft.NewInmemStore()
		rcluster.n.dir = dir
		rcluster.n.stores = store
		rcluster.n.fsm = &fsm.StorageFSM{}

		//create snap dir

	}
}
