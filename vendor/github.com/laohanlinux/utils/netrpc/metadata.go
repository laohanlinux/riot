package netrpc

import "context"

type netRpcKey struct{}
type NetRPCMetaData map[string]string

func NewMetaData(metadata map[string]string) NetRPCMetaData {
	md := NetRPCMetaData{}
	for k, v := range metadata {
		md[k] = v
	}
	return md
}

func NewContext(ctx context.Context, md NetRPCMetaData) context.Context {
	return context.WithValue(ctx, netRpcKey{}, md)
}

func FromContext(ctx context.Context) (md NetRPCMetaData, ok bool) {
	md, ok = ctx.Value(netRpcKey{}).(NetRPCMetaData)
	return md, ok
}
