package cluster

import (
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/fsm"
)

type node struct {
	dir    []string
	stores *raft.InmemStore
	fsm    *fsm.StorageFSM
	snap   *raft.FileSnapshotStore
	tran   *raft.InmemTransport
	r      *raft.Raft
}

// close node
func (c *node) Close() {
	// Wait for shutdown
	timer := time.AfterFunc(200*time.Microsecond, func() {
		panic("time out waiting for shutdown")
	})

	future := raft.Future

	if err := future.Error(); err != nil {
		panic(fmt.Errorf("shutdown future err: %v", err))
	}

	timer.Stop()

	//delete all old database
	os.RemoveAll(c.dir)
}

// func (c *cluser) GetInState(s raft.RaftState) raft.RaftState {
// 	return c.r.State()
// }

func (c *node) Leader() *raft.Raft {
	timeout := time.AfterFunc(400*time.Microsecond, func() {
		panic("timeout waitting for leader")
	})

	defer timeout.Stop()
	//            logger.Info(...)
	// if c.r.State() == raft.Leader {
	//                return
	// }
	logger.Info("No ")
	return nil
}

func (c *node) Connect() {
	logger.Info("excute in raft.Connect")

}

func (c *node) Disconnect(a string) {
	c.tran.DisconnectAll()
}
