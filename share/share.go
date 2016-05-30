package share

import (
	"encoding/json"
)

var ShCache *ShareCache

type ShareCache struct {
	LRPC             *LeaderRpcAddr `json:"lrpc"`
	LHA              *LeaderHTTAddr `json:"lha"`
	StoreBackendType string         `json:"store_backend_type"`
}

func (sc ShareCache) ToBytes() []byte {
	b, _ := json.Marshal(sc)
	return b
}

type LeaderRpcAddr struct {
	Addr string
	Port string
}

type LeaderHTTAddr struct {
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
		StoreBackendType: "",
	}
}
