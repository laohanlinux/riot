package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/cmd"
	"github.com/laohanlinux/riot/rpc/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//.
const (
	Ok = iota
	ErrRPCApply
)

func NewRpcServer(addr string) (RiotRPCService, error) {
	rrpc := RiotRPCService{}
	serverLis, err := net.Listen("tcp", addr)
	if err != nil {
		return rrpc, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRiotGossipServer(grpcServer, &rrpc)
	return rrpc, grpcServer.Serve(serverLis)
}

type RiotRPCService struct{}

// OpRPC handles rpc reuqest for set and del operation
func (rcs *RiotRPCService) OpRPC(ctx context.Context, r *pb.OpRequest) (*pb.OpReply, error) {
	var err error
	var value []byte
	switch r.Op {
	// need not to do raftNode
	case cmd.CmdGetBucket:
		var bStats interface{}
		bStats, err = cluster.SingleCluster().FSM.GetBucket([]byte(r.Bucket))
		if err == nil {
			value, err = json.Marshal(bStats)
		}
	// need not to do raftNode
	case cmd.CmdGet:
		value, err = cluster.SingleCluster().FSM.Get([]byte(r.Bucket), r.Value)
	default:
		raftNode := cluster.SingleCluster().R
		b, _ := json.Marshal(r)
		err = raftNode.Apply(b, time.Second).Error()
	}

	if err != nil {
		return &pb.OpReply{
			Status:  0,
			ErrCode: ErrRPCApply,
			Msg:     fmt.Sprintf("%v", err),
		}, nil
	}
	return &pb.OpReply{
		Status:  1,
		ErrCode: Ok,
		Msg:     "",
		Value:   value,
	}, nil
}
