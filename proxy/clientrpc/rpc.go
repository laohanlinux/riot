package clientrpc

import (
	"sync"
	"time"

	"github.com/laohanlinux/riot/api"
	log "github.com/laohanlinux/utils/gokitlog"
	"github.com/laohanlinux/utils/pool/rpc"
)

var (
	DefaultLeaderRPC *LeaderRPC
	DefaultRaftRPC   *rpc.NetRPCRing
)

const (
	APIService             = "APIService"
	APIServiceKV           = APIService + "." + "KV"
	APIServiceSetKV        = APIService + "." + "SetKV"
	APIServiceDelKey       = APIService + "." + "DelKey"
	APIServiceDelBucket    = APIService + "." + "DelBucket"
	APIServiceBucketInfo   = APIService + "." + "BucketInfo"
	APIServiceCreateBucket = APIService + "." + "CreateBucket"

	APIServiceState      = APIService + "." + "NodeState"
	APIServiceNodeString = APIService + "." + "NodeString"
	APIServicePeers      = APIService + "." + "Peers"
	APIServiceLeader     = APIService + "." + "Leader"
	APIServiceSnapshot   = APIService + "." + "Snapshot"
	APIServiceRemovePeer = APIService + "." + "RemovePeer"
)

func KV(bucketName, key string, qs int) (value []byte, has bool, err error) {
	var (
		arg   = api.GetKVArg{BucketName: bucketName, Key: key}
		reply api.GetKVReply
	)

	if qs == 1 {
		err = DefaultLeaderRPC.Call(APIServiceKV, &arg, &reply)
	} else {
		err = DefaultRaftRPC.Call(APIServiceKV, &arg, &reply)
	}
	value = reply.Value

	return
}

func SetKV(bucketName, key string, value []byte) (err error) {
	var (
		arg   = api.SetKVArg{BucketName: bucketName, Key: key, Value: value}
		reply api.NotReply
	)
	err = DefaultLeaderRPC.Call(APIServiceSetKV, &arg, &reply)
	return
}

func DelKey(bucketName, key string) (err error) {
	var (
		arg   = api.DelKVArg{BucketName: bucketName, Key: key}
		reply api.NotReply
	)

	err = DefaultLeaderRPC.Call(APIServiceDelKey, &arg, &reply)
	return
}

func DelBucket(bucketName string) (err error) {
	var (
		arg   = api.DelKVArg{BucketName: bucketName}
		reply api.NotReply
	)
	err = DefaultLeaderRPC.Call(APIServiceDelBucket, &arg, &reply)
	return
}

func BucketInfo(bucketName string) (info interface{}, has bool, err error) {
	var (
		arg   = api.BucketInfoArg{BucketName: bucketName}
		reply api.BucketInfoReply
	)
	if err = DefaultLeaderRPC.Call(APIServiceBucketInfo, &arg, &reply); err != nil {
		return
	}

	info = reply.Info
	has = reply.Has
	return
}

func CreateBucket(bucketName string) (err error) {
	var (
		arg   = api.CreateBucketArg{BucketName: bucketName}
		reply api.NotReply
	)

	err = DefaultLeaderRPC.Call(APIServiceCreateBucket, &arg, &reply)
	return
}

func States() ([]string, error) {
	var (
		arg   api.NotArg
		reply api.NodeString
		idx   = DefaultRaftRPC.Size()
		res   []string
		err   error
	)
	for ; idx > 0; idx-- {
		if err = DefaultRaftRPC.Call(APIServiceNodeString, &arg, &reply); err != nil {
			return nil, err
		}
		res = append(res, reply.NodeInfo)
	}

	return res, nil
}

func Peers() (peers []string, err error) {
	var (
		arg   api.NotArg
		reply api.PeersReply
	)

	if err = DefaultRaftRPC.Call(APIServicePeers, &arg, &reply); err != nil {
		return
	}
	return reply.Peers, nil
}

func Leader() (addr string, err error) {
	var (
		arg   api.NotArg
		reply api.LeaderReply
	)
	if err = DefaultRaftRPC.Call(APIServiceLeader, &arg, &reply); err != nil {
		return
	}
	return reply.Leader, nil
}

func Snapshot() (snaLen int, err error) {
	var (
		arg   api.NotArg
		reply api.SnapshotReply
	)
	if err = DefaultRaftRPC.Call(APIServiceSnapshot, &arg, &reply); err != nil {
		return
	}
	return reply.Len, nil
}

func RemovePeer(peer string) (err error) {
	var (
		arg   = api.RemovePeerArg{Peer: peer}
		reply api.RemovePeerReply
	)
	err = DefaultLeaderRPC.Call(APIServiceRemovePeer, &arg, &reply)
	return
}

func InitRPC(nodeAddrs []string, poolSize int) error {
	var (
		err     error
		errChan = make(chan error)
		success int
	)
	go func() {
		errChan <- initRPC(nodeAddrs, poolSize)
	}()
	go func() {
		errChan <- initLeaderRPC(nodeAddrs, poolSize)
	}()
	for {
		select {
		case <-time.After(time.Second * 3):
			panic("rpc init timeout")
		case err = <-errChan:
			if err != nil {
				return err
			}
			success++
			if success == 2 {
				return nil
			}
		}
	}
	return nil
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
		leaderRPC.opt = opt
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
	opt  rpc.NetRPCRingOpt
}

func (l *LeaderRPC) start() {
	var (
		ticker    = time.NewTicker(time.Microsecond)
		firstTime = true
		arg       api.NotArg
		reply     api.NodeStateReply
		err       error
		g         sync.WaitGroup
		nodes     = l.nodes
	)
	log.Debug("nodes", nodes)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if firstTime {
				ticker.Stop()
				ticker = time.NewTicker(time.Second * 3)
			}
			// check current node
			err = l.NetRPCRing.Call(APIServiceState, &arg, &reply)
			if err == nil {
				if reply.State == "Leader" {
					return
				}
				log.Warnf("leader change, old leader:%v, current node state:%v", l.opt.Addr, reply.State)
			}

			for _, addr := range nodes {
				g.Add(1)
				go func(addr string) {
					defer g.Done()
					var (
						opt    = rpc.NetRPCRingOpt{NetWork: "tcp", Addr: addr, PoolSize: l.opt.PoolSize}
						c, err = rpc.NewNetRPCRing([]rpc.NetRPCRingOpt{opt})
						arg    api.NotArg
						reply  api.NodeStateReply
					)
					if err != nil {
						log.Error("err", err)
						return
					}
					err = c.Call(APIServiceState, &arg, &reply)
					if err != nil {
						log.Error("err", err)
						c.Close()
						return
					}
					if reply.State == "Leader" {
						l.NetRPCRing = c
						log.Warn("new leader", addr)
					} else {
						c.Close()
						log.Warnf("the node(%v) is %v", addr, reply.State)
					}
				}(addr)
			}
			g.Wait()
		case <-l.quit:
			return
		}
	}
}

func (l *LeaderRPC) Close() {
	close(l.quit)
	l.NetRPCRing.Close()
}
