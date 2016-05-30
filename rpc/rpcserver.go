package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/laohanlinux/riot/cluster"
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

// RiotRPCService .
type RiotRPCService struct{}

// OpRPC handles rpc reuqest for set and del operation
func (rcs *RiotRPCService) OpRPC(ctx context.Context, r *pb.OpRequest) (*pb.OpReply, error) {
	b, _ := json.Marshal(r)
	// get the local fsm

	raftNode := cluster.SingleCluster().R
	future := raftNode.Apply(b, time.Second)
	if err := future.Error(); err != nil {
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
	}, nil
}
