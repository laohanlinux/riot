package command

import (
	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/pb"
)

// ....
const (
	CmdGet = "get"
	CmdSet = "set"
	CmdDel = "del"
)

type Command struct {
	Op    string
	Key   string
	Value []byte
}

// DoGet returns value by specified key
func (cm Command) DoGet() ([]byte, error) {
	c := cluster.SingleCluster()
	return c.Node().Get(cm.Key)
}

func (cm Command) DoSet() error {
	c := cluster.SingleCluster()
	addr := c.Leader().Leader()
	opRequest := pb.OpRequest{
		Op:    cm.Op,
		Key:   cm.Key,
		Value: cm.Value,
	}
	_, err := pb.NewRiotRPCClient().RPCRequest(addr, &opRequest)
	if err != nil {
		return err
	}
	return nil
}

func (cm Command) DoDel() error {
	c := cluster.SingleCluster()
	addr := c.Leader().Leader()
	opRequest := pb.OpRequest{
		Op:    cm.Op,
		Key:   cm.Key,
		Value: cm.Value,
	}
	_, err := pb.NewRiotRPCClient().RPCRequest(addr, &opRequest)
	return err
}
