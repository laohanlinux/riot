package rpc

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/laohanlinux/utils/netrpc"
)

var ErrNoAvailableRPCConn = errors.New("no available rpc connection")

func NewNetRPCRing(opts []NetRPCRingOpt) (*NetRPCRing, error) {
	var (
		ccs []*netrpcClients
		cc  *netrpcClients
		err error
	)
	if len(opts) == 0 {
		panic("invalid rpc ring config")
	}
	defer func() {
		if err != nil && ccs != nil {
			for _, cc := range ccs {
				cc.Close()
			}
		}
	}()

	for _, opt := range opts {
		cc, err = newrpcClients(opt.NetWork, opt.Addr, opt.PoolSize)
		if err != nil {
			return nil, err
		}
		fmt.Printf("create a net rpc clients, address:%v, poolSize:%d\n", opt.Addr, opt.PoolSize)
		ccs = append(ccs, cc)
	}
	return &NetRPCRing{idx: 0, pool: ccs, test: opts[0].test}, nil
}

func AddClient(netWork, hostPort string, ring *NetRPCRing) error {
	var (
		err error
		cc  *netrpcClients
		ccs []*netrpcClients
	)
	cc, err = newrpcClients(netWork, hostPort, len(ring.pool[0].sources))
	if err != nil {
		return err
	}

	ccs = append(ccs, ring.pool...)
	ccs = append(ccs, cc)
	ring.pool = ccs
	return nil
}

// TODO
// add remove function

type NetRPCRingOpt struct {
	NetWork  string
	Addr     string
	PoolSize int
	test     bool
}

type NetRPCRing struct {
	idx  uint64
	test bool
	pool []*netrpcClients
}

func (n *NetRPCRing) Call(serviceMethod string, args, reply interface{}) error {
	var (
		c   *netrpcClient
		err error
	)
	c, err = n.receiveAliceConn()
	if err != nil {
		return err
	}
	if n.test {
		fmt.Printf("target rpc server:%s\n", c.addr)
	}
	return c.Call(serviceMethod, nil, args, reply)
}

func (n *NetRPCRing) CallWithMetaData(serviceMethod string, meta map[string]string, args, reply interface{}) error {
	var (
		c   *netrpcClient
		err error
	)
	c, err = n.receiveAliceConn()
	if err != nil {
		return err
	}
	return c.Call(serviceMethod, meta, args, reply)
}

func (n *NetRPCRing) Go(serviceMethod string, done chan *netrpc.Call, args, reply interface{}) *netrpc.Call {
	var (
		c *netrpcClient
	)
	c, _ = n.receiveAliceConn()
	return c.Go(serviceMethod, nil, args, reply, done)
}

func (n *NetRPCRing) GoWithMetaData(serviceMethod string, meta map[string]string, done chan *netrpc.Call, args, reply interface{}) *netrpc.Call {
	var (
		c *netrpcClient
	)
	c, _ = n.receiveAliceConn()
	return c.Go(serviceMethod, meta, args, reply, done)
}

func (n *NetRPCRing) Size() int {
	return len(n.pool)
}

func (n *NetRPCRing) Close() {
	for _, c := range n.pool {
		c.Close()
	}
}

func (n *NetRPCRing) receiveAliceConn() (*netrpcClient, error) {
	var (
		c    *netrpcClient
		err  error
		idx  = atomic.LoadUint64(&n.idx)
		size = uint64(len(n.pool))
		i    uint64
	)
	atomic.AddUint64(&n.idx, 1)

	for ; i < size; i++ {
		j := int((idx + i) % size)
		c, err = n.pool[j].GetAliveConn()
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		break
	}

	if c == nil {
		return nil, ErrNoAvailableRPCConn
	}
	return c, nil
}

func newrpcClients(network, addr string, size int) (*netrpcClients, error) {
	var (
		ccs []*netrpcClient
		c   *netrpcClient
		err error
	)
	defer func() {
		if err != nil && ccs != nil {
			for _, cc := range ccs {
				cc.close()
			}
		}
	}()

	for i := 0; i < size; i++ {
		c, err = newNetRPCClient(network, addr)
		if err != nil {
			return nil, err
		}
		ccs = append(ccs, c)
	}

	return &netrpcClients{idx: 0, sources: ccs}, nil
}

type netrpcClients struct {
	idx     uint64
	sources []*netrpcClient
}

func (ncs *netrpcClients) GetAliveConn() (*netrpcClient, error) {
	var (
		idx  = atomic.LoadUint64(&ncs.idx)
		size = uint64(len(ncs.sources))
		i    uint64
	)
	atomic.AddUint64(&ncs.idx, 1)
	for ; i < size; i++ {
		j := int((idx + i) % size)
		if ncs.sources[j].alive {
			return ncs.sources[j], nil
		}
	}
	return nil, ErrNoAvailableRPCConn
}

func (ncs *netrpcClients) Close() {
	for _, c := range ncs.sources {
		c.close()
	}
}

func newNetRPCClient(netWork, addr string) (*netrpcClient, error) {
	nc := &netrpcClient{
		alive:   false,
		quit:    make(chan int),
		netWork: netWork,
		addr:    addr,
	}
	c, err := netrpc.Dial(netWork, addr)
	if err != nil {
		return nil, err
	}
	nc.Client = c
	nc.alive = true
	go nc.heartbeat()
	return nc, nil
}

type NetrpcClient interface {
	close()
	heartbeat()
}

type netrpcClient struct {
	alive   bool
	quit    chan int
	netWork string
	addr    string
	*netrpc.Client
}

func (nc *netrpcClient) close() {
	close(nc.quit)
}

func (nc *netrpcClient) heartbeat() {
	var (
		t          = time.NewTicker(time.Second)
		args       = netrpc.EmptyRequest{}
		reply      = netrpc.EmptyReply{}
		checkTimes = 0
		err        error
	)
	for {
		select {
		case <-t.C:
			if err = nc.Call(netrpc.HealthCheckPingNetRPC, nil, args, &reply); err != nil {
				fmt.Printf("health check error:%v\n", err.Error())
				if checkTimes > 3 {
					nc.alive = false
					checkTimes = 0
					// reconnection
					c, err := netrpc.Dial(nc.netWork, nc.addr)
					if err != nil {
						fmt.Printf("reconnection fail, err:%v\n", err)
					} else {
						//close old
						nc.Client.Close()
						nc.Client = c
					}
				}
				checkTimes++
			} else if nc.alive == false {
				nc.alive = true
				checkTimes = 0
			}
		case <-nc.quit:
			nc.alive = false
			nc.Client.Close()
			t.Stop()
			return
		}
	}
}
