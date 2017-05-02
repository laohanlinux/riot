package api

import "context"

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
	Info interface{}
}

////////////
type NodeStateReply struct {
	State string
}

//////////

type APIService struct {
	api API
	adm AdmAPI
}

func (s *APIService) KV(_ context.Context, arg *GetKVArg, reply *GetKVReply) (err error) {
	reply.Value, err = s.api.GetValue(arg.BucketName, arg.Key)
	return
}
func (s *APIService) SetKV(_ context.Context, arg *SetKVArg, _ *NotReply) (err error) {
	err = s.api.SetKV(arg.BucketName, arg.Key, arg.Value)
	return
}

func (s *APIService) BucketInfo(_ context.Context, arg *BucketInfoArg, reply *BucketInfoReply) (err error) {
	reply.Info, err = s.api.GetBucket(arg.BucketName)
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

func (s *APIService) NodeState(_ context.Context, arg *NotArg, reply *NodeStateReply) (err error) {
	reply.State = s.adm.State()
	return
}
