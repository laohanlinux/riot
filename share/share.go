package share

import (
	"encoding/json"
	"time"

	"github.com/laohanlinux/riot/cmd"
	"github.com/laohanlinux/riot/config"
	"github.com/laohanlinux/riot/rpc/pb"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
)

var ShCache *ShareCache

// ShareCache for nodes share cache object
type ShareCache struct {
	LRPC             *LeaderRpcAddr `json:"lrpc"`
	LHA              *LeaderHTTAddr `json:"lha"`
	StoreBackendType string         `json:"store_backend_type"`
}

func (sc ShareCache) ToBytes() []byte {
	b, _ := json.Marshal(sc)
	return b
}

// LeaderRpcAddr of leader node rpc address
type LeaderRpcAddr struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}

// LeaderHTTAddr of leader node http address
type LeaderHTTAddr struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}

// init the share cache content
func init() {
	ShCache = &ShareCache{
		LRPC: &LeaderRpcAddr{
			Addr: "",
			Port: "",
		},
		StoreBackendType: "",
	}
}

// UpdateShareMemory to update share memory object
func UpdateShareMemory(cfg *config.Configure, r *raft.Raft) {
	ShCache.StoreBackendType = cfg.RaftC.StoreBackend

	for {
		time.Sleep(time.Second * 3)
		// get leaderName
		if cfg.RaftC.AddrString() == r.Leader() {
			// update leader addr info
			opRequest := pb.OpRequest{
				Op:    cmd.CmdShare,
				Key:   "",
				Value: []byte(cfg.RpcC.AddrString()),
			}
			ShCache.LRPC.Addr, ShCache.LRPC.Port = cfg.RpcC.Addr, cfg.RpcC.Port
			opRequest.Value, _ = json.Marshal(ShCache)
			b, _ := json.Marshal(opRequest)
			err := r.Apply(b, 3)
			if err != nil && err.Error() != nil {
				logger.Debug(r.Leader(), err.Error())
				continue
			}
			time.Sleep(time.Second * 5)
		}
		cfg.LeaderRpcC.Addr, cfg.LeaderRpcC.Port = ShCache.LRPC.Addr, ShCache.LRPC.Port
		// all the node backend store type must be same
		if ShCache.StoreBackendType != cfg.RaftC.StoreBackend {
			logger.Fatal("all raft node backend store type muset be same")
		}
	}

}
