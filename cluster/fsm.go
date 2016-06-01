package cluster

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/laohanlinux/riot/cmd"
	"github.com/laohanlinux/riot/rpc/pb"
	"github.com/laohanlinux/riot/share"
	"github.com/laohanlinux/riot/store"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var ErrNotFound = fmt.Errorf("the key's value is nil.")
var ErrInvalidCmd = fmt.Errorf("the command is invalid")

// var ErrInvalidBackendStore = fmt.Errorf("the store backend must be boltdb")

func NewStorageFSM(rs store.RiotStorage) *StorageFSM {
	return &StorageFSM{
		l:  &sync.Mutex{},
		rs: rs,
	}
}

// StorageFSM is an implememtation of the FSM interfacec, and just
// storage the key/value logs sequentially
type StorageFSM struct {
	l  *sync.Mutex
	rs store.RiotStorage
}

// Apply is noly call in out with master leader
// log format: json
// {"cmd":op, "key":key, "value": value}
// TODO
// use protocol buffer instead of json format
func (s *StorageFSM) Apply(log *raft.Log) interface{} {
	s.l.Lock()
	defer s.l.Unlock()
	//logger.Debug("Excute StorageFSM.Apply ...")
	var req pb.OpRequest
	if err := json.Unmarshal(log.Data, &req); err != nil {
		logger.Fatal(err)
	}

	var err error
	switch req.Op {
	case cmd.CmdSet:
		err = s.rs.Set([]byte(req.Bucket), []byte(req.Key), req.Value)
	case cmd.CmdDel:
		err = s.rs.Del([]byte(req.Bucket), []byte(req.Key))
	case cmd.CmdCreateBucket:
		rs, _ := s.rs.(*store.BoltdbStore)
		err = rs.CreateBucket([]byte(req.Bucket))
	case cmd.CmdDelBucket:
		rs, _ := s.rs.(*store.BoltdbStore)
		err = rs.DelBucket([]byte(req.Bucket))
	case cmd.CmdShare:
		err = json.Unmarshal(req.Value, share.ShCache)
	default:
		err = ErrInvalidCmd
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// Get .
func (s *StorageFSM) Get(bucket, key []byte) ([]byte, error) {
	s.l.Lock()
	defer s.l.Unlock()
	value, err := s.rs.Get(bucket, key)
	if err == errors.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *StorageFSM) GetBucket(bucket []byte) (interface{}, error) {
	s.l.Lock()
	defer s.l.Unlock()
	rs, _ := s.rs.(*store.BoltdbStore)
	return rs.GetBucket([]byte(bucket))
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
	//logger.Info("Must clear old dirty data, Excute StorageFSN.Restore ...")
	s.l.Lock()
	defer s.l.Unlock()
	defer inp.Close()

	bSizeBuf := make([]byte, 2)
	iterm := store.Iterm{}
	for {
		_, err := inp.Read(bSizeBuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		bSize := int(binary.LittleEndian.Uint16(bSizeBuf))
		buf := make([]byte, bSize)
		_, err = inp.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic("snapshot decode error:" + err.Error())
		}
		// decoding
		if err = json.Unmarshal(buf, &iterm); err != nil {
			errMsg := fmt.Sprintf("decode json(%s) error in restore snapshot, error:%s", buf, err.Error())
			panic(errMsg)
		}
		if err = s.rs.Set(iterm.Bucket, iterm.Key, iterm.Value); err != nil {
			panic("restore data into backend store happends error: " + err.Error())
		}
	}

	return nil
}

// StorageSnapshot .
type StorageSnapshot struct {
	diskStore store.RiotStorage
}

// Persist ...
func (s *StorageSnapshot) Persist(sink raft.SnapshotSink) error {
	logger.Info("Excute StorageSnapshot.Persist ... ")
	defer sink.Close()
	c := s.diskStore.Rec()

	for {
		iterm := <-c
		if iterm.Err == nil {
			data, err := json.Marshal(iterm)
			if err != nil {
				return err
			}
			bSize := uint16(len(data))
			buf := make([]byte, bSize+2)
			binary.LittleEndian.PutUint16(buf[:2], bSize)
			copy(buf[2:], data)
			if _, err = sink.Write(buf); err != nil {
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
	return conf
}
