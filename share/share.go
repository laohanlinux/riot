package share

import "encoding/json"

var ShCache *ShareCache

type ShareCache struct {
	LRPC *LeaderRpcAddr
}

func (sc ShareCache) ToBytes() []byte {
	b, _ := json.Marshal(sc)
	return b
}

type LeaderRpcAddr struct {
	Addr string
	Port string
}

// init the share cache content
func init() {
	ShCache = &ShareCache{
		LRPC: &LeaderRpcAddr{
			Addr: "",
			Port: "",
		},
	}
}
