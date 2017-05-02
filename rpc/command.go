package rpc

// RpcCmd of rpc remote call
type RpcCmd struct {
	Op     string
	Bucket string
	Key    string
	Value  []byte
}
