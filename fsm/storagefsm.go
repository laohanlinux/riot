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
	"github.com/laohanlinux/riot/pb"
)

// StorageFSM is an implememtation of the FSM interfacec, and just
// storage the key/value logs sequentially
type StorageFSM struct {
	sync.Mutex
	cache map[string][]byte
}

// Apply is noly call in out with master leader
// log format: json
// {"cmd":op, "key":key, "value": value}
// TODO
// use protocol buffer instead of json format
func (s *StorageFSM) Apply(log *raft.Log) interface{} {
	s.Lock()
	defer s.Unlock()

	logger.Info("Excute StorageFSM.Apply ...")
	var req pb.OpRequest
	if err := json.Unmarshal(log.Data, &req); err != nil {
		logger.Fatal(err)
	}

	switch req.Op {
	case "Get":
		s.cache[req.Key] = req.Value
	case "SET":
		delete(s.cache, req.Key)
	default:
		return fmt.Errorf("%s is a invalid command", req.Op)
	}

	return nil
}

// Get .
func (s *StorageFSM) Get(key string) ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	value, ok := s.cache[key]
	if !ok {
		return nil, fmt.Errorf("not found the %s's value\n", key)
	}
	return value, nil
}

// Snapshot .
func (s *StorageFSM) Snapshot() (raft.FSMSnapshot, error) {
	s.Lock()
	defer s.Unlock()
	logger.Info("Excute StorageFSM.Snapshot ...")
	// return &StorageSnapshot{s.logs, len(s.logs)}, nil
	return &StorageSnapshot{
		diskCache: s.cache,
	}, nil
}

// Restore data from persit location
func (s *StorageFSM) Restore(inp io.ReadCloser) error {
	logger.Info("Excute StorageFSN.Restore ...")
	s.Lock()
	defer s.Unlock()
	defer inp.Close()
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(inp, &hd)
	s.cache = nil
	return dec.Decode(&s.cache)
}

// StorageSnapshot .
type StorageSnapshot struct {
	// logs     [][]byte
	// maxIndex int
	diskCache map[string][]byte
}

// Persist ...
func (s *StorageSnapshot) Persist(sink raft.SnapshotSink) error {
	logger.Info("Excute StorageSnapshot.Persist ... ")

	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(sink, &hd)

	if err := enc.Encode(s.diskCache); err != nil {
		sink.Close()
		return err
	}
	sink.Close()
	return nil
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
