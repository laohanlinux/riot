package netrpc

import "golang.org/x/net/context"

const (
	HealthCheckService    = "HealthCheck"
	HealthCheckPingNetRPC = "HealthCheck.Ping"
)

type HealthCheck struct{}

func (hc *HealthCheck) Ping(_ context.Context, req *EmptyRequest, reply *EmptyReply) error {
	return nil
}
