package command

import (
	//"fmt"
	"fmt"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/rpc"
	"github.com/laohanlinux/riot/rpc/pb"
)

// ....
const (
	CmdGet   = "GET"
	CmdSet   = "SET"
	CmdDel   = "DEL"
	CmdShare = "SHARE"
)

type Command struct {
	Op    string
	Key   string
	Value []byte
}

// DoGet returns value by specified key
// TODO
// search a value from Leader Node or Follower Node mannully
func (cm Command) DoGet() ([]byte, error) {
	c := cluster.SingleCluster()
	return c.Get(cm.Key)
}

func (cm Command) DoSet() error {
	cfg := config.GetConfigure()
	// if cfg.LeaderRpcC.Addr == "" {
	// 	//update Leader Info
	// 	c := cluster.SingleCluster()
	// 	addr := strings.Split(c.Leader(), ":")
	// 	cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port = addr[0], addr[1]
	// }
	rpcAddr := fmt.Sprintf("%s:%s", cfg.LeaderRpcC.Addr, cfg.RpcC.Port)
	opRequest := pb.OpRequest{
		Op:    cm.Op,
		Key:   cm.Key,
		Value: cm.Value,
	}
	reply, err := rpc.NewRiotRPCClient().RPCRequest(rpcAddr, &opRequest)
	if reply.ErrCode != 1 {
		err = fmt.Errorf("%s", reply.Msg)
	}
	return err
}

func (cm Command) DoDel() error {
	cfg := config.GetConfigure()
	// if cfg.LeaderRpcC.Addr == "" {
	// 	//update Leader Info
	// 	c := cluster.SingleCluster()
	// 	addr := strings.Split(c.Leader(), ":")
	// 	cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port = addr[0], addr[1]
	// }
	rpcAddr := fmt.Sprintf("%s:%s", cfg.LeaderRpcC.Addr, cfg.RpcC.Port)
	opRequest := pb.OpRequest{
		Op:    cm.Op,
		Key:   cm.Key,
		Value: cm.Value,
	}
	reply, err := rpc.NewRiotRPCClient().RPCRequest(rpcAddr, &opRequest)
	if reply.Status != 1 {
		err = fmt.Errorf("%s", reply.Msg)
	}
	return err
}
