package rpc

import (
	"fmt"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/cmd"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/rpc/pb"
)

type RpcCmd struct {
	Op     string
	Bucket string
	Key    string
	Value  []byte
}

// DoGet returns value by specified key
func (rc RpcCmd) DoGet(qs int) ([]byte, error) {
	switch qs {
	case cmd.QsConsistent:
		// make it simple better, proxy to leader http port
		cfg := config.GetConfigure()
		rpcAddr := fmt.Sprintf("%s:%s", cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port)
		opRequest := pb.OpRequest{
			Op:     rc.Op,
			Bucket: rc.Bucket,
			Key:    rc.Key,
			Value:  rc.Value,
		}
		reply, err := NewRiotRPCClient().RPCRequest(rpcAddr, &opRequest)
		if err != nil {
			return nil, err
		}
		if reply.Status != 1 {
			return nil, fmt.Errorf("%s", reply.Msg)
		}
		return reply.Value, nil
	case cmd.QsRandom:
		c := cluster.SingleCluster()
		return c.Get([]byte(rc.Bucket), []byte(rc.Key))
	default:
		return nil, fmt.Errorf("the qury strategies is invalid.")
	}
}

func (rc RpcCmd) DoSet() error {
	cfg := config.GetConfigure()
	rpcAddr := fmt.Sprintf("%s:%s", cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port)
	opRequest := pb.OpRequest{
		Op:     rc.Op,
		Bucket: rc.Bucket,
		Key:    rc.Key,
		Value:  rc.Value,
	}
	reply, err := NewRiotRPCClient().RPCRequest(rpcAddr, &opRequest)
	if reply.Status != 1 {
		err = fmt.Errorf("%s", reply.Msg)
	}
	return err
}

func (rc RpcCmd) DoDel() error {
	cfg := config.GetConfigure()
	rpcAddr := fmt.Sprintf("%s:%s", cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port)
	opRequest := pb.OpRequest{
		Op:     rc.Op,
		Bucket: rc.Bucket,
		Key:    rc.Key,
		Value:  rc.Value,
	}
	reply, err := NewRiotRPCClient().RPCRequest(rpcAddr, &opRequest)
	if reply.Status != 1 {
		err = fmt.Errorf("%s", reply.Msg)
	}
	return err
}

func (rc RpcCmd) GetBucket() (interface{}, error) {
	cfg := config.GetConfigure()
	rpcAddr := fmt.Sprintf("%s:%s", cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port)
	opRequest := pb.OpRequest{
		Op:     rc.Op,
		Bucket: rc.Bucket,
		Key:    rc.Key,
		Value:  rc.Value,
	}
	reply, err := NewRiotRPCClient().RPCRequest(rpcAddr, &opRequest)
	if err != nil {
		return nil, err
	}
	if reply.Status != 1 {
		return nil, fmt.Errorf("%s", reply.Msg)
	}
	return reply.Value, nil
}
