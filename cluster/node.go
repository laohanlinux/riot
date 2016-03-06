package cluster

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/fsm"
)

type Node struct {
	addr   string
	dir    string
	stores *raft.InmemStore
	fsm    *fsm.StorageFSM
	snap   *raft.FileSnapshotStore
	tran   *raft.NetworkTransport
	r      *raft.Raft
}

// close node
func (c *Node) Close() {
	// Wait for shutdown
	timer := time.AfterFunc(200*time.Microsecond, func() {
		panic("time out waiting for shutdown")
	})

	var future raft.Future

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

func (c *Node) Leader() *raft.Raft {
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

func (c *Node) GetFSM() raft.FSM {
	return c.fsm
}

func (c *Node) Connect() {
	logger.Info("excute in raft.Connect")

}

func (c *Node) Disconnect(a string) {
	c.tran.Close()
}

// Get .
func (n *Node) Get(Key string) ([]byte, error) {
	return n.fsm.Get(Key)
}
