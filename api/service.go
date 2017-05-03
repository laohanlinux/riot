package api

import (
	"context"

	"github.com/boltdb/bolt"
	"github.com/laohanlinux/riot/cluster"
)

type NotArg struct{}
type NotReply struct{}

/////////////////////////
type SetKVArg struct {
	BucketName string
	Key        string
	Value      []byte
}

type GetKVArg struct {
	BucketName string
	Key        string
}

type GetKVReply struct {
	Has   bool
	Value []byte
}

type DelKVArg struct {
	BucketName string
	Key        string
}

//////////////

type DelBucketArg struct {
	BucketName string
}

type BucketInfoArg struct {
	BucketName string
}

type CreateBucketArg struct {
	BucketName string
}

type BucketInfoReply struct {
	Has  bool
	Info bolt.BucketStats
}

////////////
type NodeStateReply struct {
	State string
}

type PeersReply struct {
	Peers []string
}

type LeaderReply struct {
	Leader string
}

type SnapshotReply struct {
	Len int
}

type RemovePeerArg struct {
	Peer string
}

//////////

type APIService struct {
	api API
	adm AdmAPI
}

func NewAPIService(api API, adm AdmAPI) *APIService {
	return &APIService{api: api, adm: adm}
}

func (s *APIService) KV(_ context.Context, arg *GetKVArg, reply *GetKVReply) (err error) {
	reply.Value, err = s.api.GetValue(arg.BucketName, arg.Key)
	if err == cluster.ErrNotFound {
		err = nil
	}
	reply.Has = true
	return
}
func (s *APIService) SetKV(_ context.Context, arg *SetKVArg, _ *NotReply) (err error) {
	err = s.api.SetKV(arg.BucketName, arg.Key, arg.Value)
	return
}

func (s *APIService) BucketInfo(_ context.Context, arg *BucketInfoArg, reply *BucketInfoReply) (err error) {
	var info interface{}
	info, err = s.api.GetBucket(arg.BucketName)
	if err == bolt.ErrBucketNotFound {
		err = nil
		return
	}
	reply.Info = info.(bolt.BucketStats)
	reply.Has = true
	return
}

func (s *APIService) DelKey(_ context.Context, arg *DelKVArg, _ *NotReply) (err error) {
	err = s.api.DelKey(arg.BucketName, arg.Key)
	return
}

func (s *APIService) DelBucket(_ context.Context, arg *DelBucketArg, _ *NotReply) (err error) {
	err = s.api.DelBucket(arg.BucketName)
	return
}

func (s *APIService) CreateBucket(_ context.Context, arg *CreateBucketArg, _ *NotReply) (err error) {
	err = s.api.CreateBucket(arg.BucketName)
	return
}

func (s *APIService) NodeState(_ context.Context, _ *NotArg, reply *NodeStateReply) (err error) {
	reply.State = s.adm.State()
	return
}

func (s *APIService) Peers(_ context.Context, _ *NotArg, reply *PeersReply) (err error) {
	reply.Peers, err = s.adm.Peers()
	return
}

func (s *APIService) Leader(_ context.Context, _ *NotArg, reply *LeaderReply) (err error) {
	reply.Leader, err = s.adm.Leader()
	return
}

func (s *APIService) Snapshot(_ context.Context, _ *NotArg, reply *SnapshotReply) (err error) {
	reply.Len, err = s.adm.Snapshot()
	return
}

func (s *APIService) RemovePeer(_ context.Context, arg *RemovePeerArg, _ *NotReply) (err error) {
	err = s.adm.RemovePeer(arg.Peer)
	return
}
