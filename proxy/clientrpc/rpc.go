package clientrpc

import (
	"sync"
	"time"

	"github.com/laohanlinux/riot/api"
	"github.com/laohanlinux/utils/pool/rpc"
	"github.com/lunny/log"
)

const (
	APIService             = "APIService"
	APIServiceKV           = APIService + "." + "KV"
	APIServiceSetKV        = APIService + "." + "SetKV"
	APIServiceDel          = APIService + "." + "Del"
	APIServiceDelKey       = APIService + "." + "DelKey"
	APIServiceDelBucket    = APIService + "." + "DelBucket"
	APIServiceBucketInfo   = APIService + "." + "BucketInfo"
	APIServiceCreateBucket = APIService + "." + "CreateBucket"
	APIServiceState        = APIService + "." + "NodeState"
)

var DefaultLeaderRPC *LeaderRPC
var DefaultRaftRPC *rpc.NetRPCRing

func InitRPC(nodeAddr []string, poolSize int) error {
	var err error
	if err = initRPC(nodeAddr, poolSize); err != nil {
		return err
	}
	if err = initLeaderRPC(nodeAdds, poolSize); err != nil {
		return err
	}
}

func initRPC(nodeAddr []string, poolSize int) error {
	var (
		opts []rpc.NetRPCRingOpt
		err  error
	)
	for _, addr := range nodeAddr {
		opts = append(opts, rpc.NetRPCRingOpt{Addr: addr, PoolSize: poolSize, NetWork: "tcp"})
	}
	DefaultRaftRPC, err = rpc.NewNetRPCRing(opts)
	return err
}

func initLeaderRPC(nodeAdds []string, poolSize int) error {
	var (
		addr = nodeAdds[0]
		err  error
	)
	DefaultLeaderRPC, err = NewLeaderRPC(nodeAdds, addr, poolSize)
	return err
}

func NewLeaderRPC(nodeAdds []string, addr string, poolSize int) (leaderRPC *LeaderRPC, err error) {
	var (
		opt = rpc.NetRPCRingOpt{NetWork: "tcp", Addr: addr, PoolSize: poolSize}
	)
	leaderRPC = &LeaderRPC{}
	{
		leaderRPC.nodes = nodeAdds
		leaderRPC.quit = make(chan struct{})
		leaderRPC.NetRPCRing, err = rpc.NewNetRPCRing([]rpc.NetRPCRingOpt{opt})
		leaderRPC.opts = opt
	}
	if err != nil {
		return
	}
	go leaderRPC.start()
	return
}

type LeaderRPC struct {
	nodes []string
	*rpc.NetRPCRing
	quit chan struct{}
	opts rpc.NetRPCRingOpt
}

func (l *LeaderRPC) start() {
	var (
		ticker    = time.NewTicker(time.Microsecond)
		firstTime = true
		arg       api.NotArg
		reply     api.NodeStateReply
		err       error
		g         sync.WaitGroup
	)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if firstTime {
				ticker.Stop()
				ticker = time.NewTicker(time.Second * 3)
			}
			// check current node
			err = l.NetRPCRing.Call(APIServiceState, &args, &reply)
			if err == nil {
				return
			}

			for _, addr := range l.nodes {
				g.Add(1)
				go func() {
					defer g.Done()
					var (
						c, err = rpc.NewNetRPCRing([]rpc.NetRPCRingOpt{opt})
						arg    api.NotArg
						reply  api.NodeStateReply
					)
					if err != nil {
						log.Error("err", err)
						return
					}
					err = l.NetRPCRing.Call(APIServiceState, &args, &reply)
					if err == nil {
						log.Error("err", err)
						return
					}
					l.NetRPCRing = c
				}()
			}
			g.Wait()
		case <-l.quit:
			return
		}
	}
}

func (l *LeaderRPC) Close() {
	close(l.quit)
}
