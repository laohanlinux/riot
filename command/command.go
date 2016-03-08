package command

import (
	"fmt"
	"strings"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/rpc"
	"github.com/laohanlinux/riot/rpc/pb"
)

// ....
const (
	CmdGet = "GET"
	CmdSet = "SET"
	CmdDel = "DEL"
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
	addr := strings.Split(c.Leader(), ":")
	fmt.Println(addr)
	cfg := config.GetConfigure()
	rpcAddr := addr[0] + ":" + cfg.RpcC.Port
	opRequest := pb.OpRequest{
		Op:    cm.Op,
		Key:   cm.Key,
		Value: cm.Value,
	}
	_, err := rpc.NewRiotRPCClient().RPCRequest(rpcAddr, &opRequest)
	if err != nil {
		return err
	}

	return nil
}

func (cm Command) DoDel() error {
	c := cluster.SingleCluster()
	addr := c.Leader()
	opRequest := pb.OpRequest{
		Op:    cm.Op,
		Key:   cm.Key,
		Value: cm.Value,
	}
	_, err := rpc.NewRiotRPCClient().RPCRequest(addr, &opRequest)
	return err
}
