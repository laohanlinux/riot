package fsm

import (
	"io"
	"sync"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
)

// StorageFSM is an implememtation of the FSM interfacec, and just
// storage the key/value logs sequentially

type StorageFSM struct {
	sync.Mutex
	logs [][]byte
}

func (s *StorageFSM) Apply(log *raft.Log) interface{} {
	s.Lock()
	defer s.Unlock()

	logger.Info("Excute StorageFSM.Apply ...")
	s.logs = append(s.logs, log.Data)
	return len(s.logs)
}

func (s *StorageFSM) Snapshot() (raft.FSMSnapshot, error) {
	s.Lock()
	defer s.Unlock()
	logger.Info("Excute StorageFSM.Snapshot ...")
	return &StorageSnapshot{s.logs, len(s.logs)}, nil
}

// restore data from persit location
func (s *StorageFSM) Restore(inp io.ReadCloser) error {
	logger.Info("Excute StorageFSN.Restore ...")
	s.Lock()
	defer s.Unlock()
	defer inp.Close()

	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(inp, &hd)

	s.logs = nil
	return dec.Decode(&s.logs)
}

type StorageSnapshot struct {
	logs     [][]byte
	maxIndex int
}
