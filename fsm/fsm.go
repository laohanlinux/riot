package fsm

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/laohanlinux/riot/rpc/pb"
	"github.com/laohanlinux/riot/store"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var ErrNotFound = fmt.Errorf("the key's value is nil.")

//
// const (
// 	AppendAction = "rebuild the storage."
// 	CreateAction = "create a new db if the db is not exist"
// 	TruncAction  = "create a new db"
// )

func NewStorageFSM(rs RiotStorage) *StorageFSM {
	return &StorageFSM{
		l:     &sync.Mutex{},
		cache: make(map[string][]byte),
		rs:    rs,
	}
}

// StorageFSM is an implememtation of the FSM interfacec, and just
// storage the key/value logs sequentially
type StorageFSM struct {
	l     *sync.Mutex
	cache map[string][]byte
	rs    RiotStorage
}

// Apply is noly call in out with master leader
// log format: json
// {"cmd":op, "key":key, "value": value}
// TODO
// use protocol buffer instead of json format
func (s *StorageFSM) Apply(log *raft.Log) interface{} {
	s.l.Lock()
	defer s.l.Unlock()

	logger.Info("Excute StorageFSM.Apply ...")
	var req pb.OpRequest
	if err := json.Unmarshal(log.Data, &req); err != nil {
		logger.Fatal(err)
	}

	var err error
	switch req.Op {
	case "SET":
		logger.Info("Set:", req.Key, req.Value)
		err = s.rs.Set([]byte(req.Key), req.Value)
	case "DEL":
		err = s.rs.Del([]byte(req.Key))
	default:
		return fmt.Errorf("%s is a invalid command", req.Op)
	}

	return err
}

// Get .
func (s *StorageFSM) Get(key string) ([]byte, error) {
	s.l.Lock()
	defer s.l.Unlock()
	value, err := s.rs.Get([]byte(key))
	logger.Info("Get:", key)
	if err == errors.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Snapshot .
func (s *StorageFSM) Snapshot() (raft.FSMSnapshot, error) {
	s.l.Lock()
	defer s.l.Unlock()
	logger.Info("Excute StorageFSM.Snapshot ...")
	return &StorageSnapshot{
		diskStore: s.rs,
	}, nil
}

// Restore data from persit location
func (s *StorageFSM) Restore(inp io.ReadCloser) error {
	logger.Info("Excute StorageFSN.Restore ...")
	s.l.Lock()
	defer s.l.Unlock()
	defer inp.Close()
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(inp, &hd)
	s.cache = nil

	return dec.Decode(&s.cache)
}

// StorageSnapshot .
type StorageSnapshot struct {
	diskCache map[string][]byte
	diskStore RiotStorage
}

// Persist ...
func (s *StorageSnapshot) Persist(sink raft.SnapshotSink) error {
	logger.Info("Excute StorageSnapshot.Persist ... ")
	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(sink, &hd)
	c := s.diskStore.Rec()
	defer sink.Close()
	for {
		iterm := <-c
		if iterm.Err == nil {
			if err := enc.Encode(iterm); err != nil {
				return err
			}
		}
		if iterm.Err == store.ErrFinished {
			return nil
		}
	}
}

// Release .
func (s *StorageSnapshot) Release() {
	logger.Info("Excute StorageSnapshot.Release ...")
}

//InmemConfig .
//configurations optimized for in-memeory
func InmemConfig() *raft.Config {
	conf := raft.DefaultConfig()
	conf.HeartbeatTimeout = 50 * time.Millisecond
	conf.ElectionTimeout = 50 * time.Millisecond
	conf.LeaderLeaseTimeout = 50 * time.Millisecond
	conf.CommitTimeout = time.Millisecond
	conf.EnableSingleNode = true
	return conf
}
