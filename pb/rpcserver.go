package pb

import "golang.org/x/net/context"

//.
const (
	Ok = iota
	ErrRPCApply
)

// RiotRPCService .
type RiotRPCService struct{}

// OpRPC handles rpc reuqest for set and del operation
func (rcs *RiotRPCService) OpRPC(ctx *context.Context, r *OpRequest) (*OpReply, error) {
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
