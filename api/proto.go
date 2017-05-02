package api

import (
	"encoding/json"
	"time"

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

type AdmAPI interface {
	State() string
}

type miniAPI struct {
	c *cluster.Cluster
}

func (api *miniAPI) GetValue(bucketName, key string) (value []byte, err error) {
	if value, err = api.c.Get([]byte(bucketName), []byte(key)); err != nil {
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
		req      = cluster.OpRequest{Op: cmd.CmdCreateBucket}
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

func (api *miniAPI) GetBucket(bucektName string) (info interface{}, err error) {
	if info, err = api.c.FSM.GetBucket([]byte(bucektName)); err != nil {
		log.Error("err", err)
	}
	return
}

type admAPI struct {
	c *cluster.Cluster
}

func (adm *admAPI) State() string {
	return cluster.SingleCluster().Status()
}
