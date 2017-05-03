package api

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/riot/cluster"
	"github.com/laohanlinux/riot/cmd"

	log "github.com/laohanlinux/utils/gokitlog"
)

type API interface {
	GetValue(bucektName, key string) (value []byte, err error)
	DelKey(bucketName, key string) (err error)
	SetKV(bucketName, key string, value []byte) (err error)
	CreateBucket(bucketName string) (err error)
	DelBucket(bucketName string) (err error)
	GetBucket(bucketName string) (info interface{}, err error)
}

type miniAPI struct {
	c *cluster.Cluster
}

func NewMiniAPI(c *cluster.Cluster) API {
	api := miniAPI{c: c}
	return &api
}

func (api *miniAPI) GetValue(bucketName, key string) (value []byte, err error) {
	if value, err = api.c.Get([]byte(bucketName), []byte(key)); err != nil && err != cluster.ErrNotFound {
		log.Error("err", err)
	}
	return
}

func (api *miniAPI) DelKey(bucektName, key string) (err error) {
	var (
		req      = cluster.OpRequest{Op: cmd.CmdDel, Key: key, Bucket: bucektName}
		b        []byte
		raftNode = cluster.SingleCluster().R
	)
	b, err = json.Marshal(req)
	if err != nil {
		return
	}
	err = raftNode.Apply(b, time.Second).Error()
	return
}

func (api *miniAPI) SetKV(bucketName, key string, value []byte) (err error) {
	var (
		req      = cluster.OpRequest{Op: cmd.CmdSet, Bucket: bucketName, Key: key, Value: value}
		b        []byte
		raftNode = cluster.SingleCluster().R
	)
	b, err = json.Marshal(req)
	if err != nil {
		return
	}
	err = raftNode.Apply(b, time.Second).Error()
	return
}

func (api *miniAPI) CreateBucket(bucketName string) (err error) {
	var (
		req      = cluster.OpRequest{Op: cmd.CmdCreateBucket, Bucket: bucketName}
		b        []byte
		raftNode = cluster.SingleCluster().R
	)
	b, err = json.Marshal(req)
	if err != nil {
		return
	}
	err = raftNode.Apply(b, time.Second).Error()
	return
}

func (api *miniAPI) DelBucket(bucketName string) (err error) {
	var (
		req      = cluster.OpRequest{Op: cmd.CmdDelBucket, Bucket: bucketName}
		b        []byte
		raftNode = cluster.SingleCluster().R
	)
	b, err = json.Marshal(req)
	if err != nil {
		return
	}
	err = raftNode.Apply(b, time.Second).Error()
	return
}

func (api *miniAPI) GetBucket(bucketName string) (info interface{}, err error) {
	if info, err = api.c.FSM.GetBucket([]byte(bucketName)); err != nil {
		log.Error("err", err)
	}
	return
}

type AdmAPI interface {
	State() string
	Peers() (peers []string, err error)
	Leader() (node string, err error)
	Snapshot() (int, error)
	RemovePeer(peer string) (err error)
}

type admAPI struct {
	c *cluster.Cluster
}

func NewAdmAPI(c *cluster.Cluster) AdmAPI {
	adm := &admAPI{c: c}
	return adm
}

func (adm *admAPI) State() string {
	return cluster.SingleCluster().Status()
}

func (adm *admAPI) Peers() (peers []string, err error) {
	peers, err = cluster.SingleCluster().PeerStorage.Peers()
	return
}

func (adm *admAPI) Leader() (node string, err error) {
	node = cluster.SingleCluster().Leader()
	if node == "" {
		err = raft.ErrNotLeader
	}
	return
}

// TODO
// only use int leader node
func (adm *admAPI) Snapshot() (snapLen int, err error) {
	var (
		sna       = cluster.SingleCluster().Snap
		snaFuture = cluster.SingleCluster().R.Snapshot()
		snapMetas []*raft.SnapshotMeta
	)

	if err = snaFuture.Error(); err != nil {
		return
	}
	if snapMetas, err = sna.List(); err != nil {
		return
	}

	return len(snapMetas), nil
}

func (adm *admAPI) RemovePeer(peer string) (err error) {
	return cluster.SingleCluster().R.RemovePeer(peer).Error()
}
