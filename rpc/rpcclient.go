package rpc

import (
	"sync"

	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/rpc/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var rc *RiotRPCClient

type RiotRPCClient struct {
	l    *sync.RWMutex
	conn map[string]*grpc.ClientConn
}

// not thread safely
func NewRiotRPCClient() *RiotRPCClient {
	if rc == nil {
		rc = &RiotRPCClient{
			l:    &sync.RWMutex{},
			conn: make(map[string]*grpc.ClientConn),
		}
	}
	return rc
}

func (rc *RiotRPCClient) RPCRequest(rpcAddr string, r *pb.OpRequest) (*pb.OpReply, error) {
	rc.l.Lock()
	var err error
	conn, ok := rc.conn[rpcAddr]
	if !ok {
		logger.Warn("New RPC Connect...")
		conn, err = grpc.Dial(rpcAddr, grpc.WithInsecure())
		if err != nil {
			rc.l.Unlock()
			return nil, err
		}
		rc.conn[rpcAddr] = conn
	}
	rc.l.Unlock()
	client := pb.NewRiotGossipClient(conn)
	return client.OpRPC(context.Background(), r)
}
