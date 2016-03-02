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
	"github.com/laohanlinux/riot/command"
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
	var cmdData command.Comand
	if err := json.Unmarshal(log.Data, &cmdData); err != nil {
		logger.Fatal(err)
	}

	switch cmdData.Op {
	case command.CmdSet:
		s.cache[cmdData.Key] = cmdData.Value
	case command.CmdDel:
		delete(s.cache, cmdData.Key)
	default:
		return fmt.Errorf("%s is a invalid command", cmdData.Op)
	}

	return nil
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

	// hd := codec.MsgpackHandle{}
	// dec := codec.NewDecoder(inp, &hd)

	// s.logs = nil
	// return dec.Decode(&s.logs)

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

	// hd := codec.MsgpackHandle{}
	// enc := codec.NewEncoder(sink, &hd)
	// logger.Info(len(s.logs))

	// if err := enc.Encode(s.logs[:s.maxIndex]); err != nil {
	// 	sink.Close()
	// 	return err
	// }

	// logger.Info(len(s.logs))
	// sink.Close()
	// return nil

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
func InmemConfig() *raft.Log {
	conf := raft.DefaultConfig()
	conf.HeartbeatTimeout = 50 * time.Millisecond
	conf.ElectionTimeout = 50 * time.Millisecond
	conf.LeaderLeaseTimeout = 50 * time.Millisecond
	conf.CommitTimeout = time.Millisecond
	conf.EnableSingleNode = true
	return conf
}
