package pb

import (
	"net"

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
	RegisterRiotGossipServer(grpcServer, &rrpc)
	return rrpc, grpcServer.Serve(serverLis)
}

// RiotRPCService .
type RiotRPCService struct{}

// OpRPC handles rpc reuqest for set and del operation
func (rcs *RiotRPCService) OpRPC(ctx context.Context, r *OpRequest) (*OpReply, error) {
	// b, _ := json.Marshal(r)

	// fsm := cluster.SingleCluster().FSM
	// if err := fsm.Apply(b); err != nil {
	// 	return &OpReply{
	// 		Status:  0,
	// 		ErrCode: ErrRPCApply,
	// 		Msg:     fmt.Sprintf("%v", err),
	// 	}, nil
	// }

	// return &OpReply{
	// 	Status:  1,
	// 	ErrCode: Ok,
	// 	Msg:     "",
	// }, nil
	return nil, nil
}
