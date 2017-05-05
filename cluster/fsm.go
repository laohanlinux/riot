package cluster

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/laohanlinux/riot/cmd"
	"github.com/laohanlinux/riot/store"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/raft"
	log "github.com/laohanlinux/utils/gokitlog"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var ErrNotFound = fmt.Errorf("the key's value is nil.")
var ErrInvalidCmd = fmt.Errorf("the command is invalid")

type OpRequest struct {
	Op     string
	Key    string
	Bucket string
	Value  []byte
}

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
func (s *StorageFSM) Apply(logEntry *raft.Log) interface{} {
	s.l.Lock()
	defer s.l.Unlock()

	var (
		req OpRequest
		err error
	)
	if err = json.Unmarshal(logEntry.Data, &req); err != nil {
		log.Crit("err", err)
	}

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
	default:
		err = ErrInvalidCmd
	}
	if err == bolt.ErrBucketNotFound {
		err = ErrNotFound
	}
	if err == bolt.ErrBucketExists {
		err = nil
	}
	if err != nil && err != ErrNotFound {
		log.Error(err.Error())
	}
	return err
}

// Get a value by bucketName and key
func (s *StorageFSM) Get(bucket, key []byte) ([]byte, error) {
	s.l.Lock()
	defer s.l.Unlock()

	var (
		value []byte
		err   error
	)
	value, err = s.rs.Get(bucket, key)
	if err == bolt.ErrBucketNotFound || err == errors.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

// GetBucket return the bucket detail info
func (s *StorageFSM) GetBucket(bucket []byte) (info interface{}, err error) {
	s.l.Lock()
	defer s.l.Unlock()
	var rs *store.BoltdbStore
	rs, _ = s.rs.(*store.BoltdbStore)
	info, err = rs.GetBucket(bucket)
	if err == bolt.ErrBucketNotFound {
		err = ErrNotFound
	}
	return
}

// Snapshot fsm statation
func (s *StorageFSM) Snapshot() (raft.FSMSnapshot, error) {
	s.l.Lock()
	defer s.l.Unlock()
	log.Info("Excute StorageFSM.Snapshot ...")
	return &StorageSnapshot{
		diskStore: s.rs,
	}, nil
}

// Restore data from disk
func (s *StorageFSM) Restore(inp io.ReadCloser) error {
	//logger.Info("Must clear old dirty data, Excute StorageFSN.Restore ...")
	s.l.Lock()
	defer s.l.Unlock()
	defer inp.Close()
	var (
		bSizeBuf = make([]byte, 2)
		iterm    store.Iterm
		bSize    int
		buf      []byte
		rs       *store.BoltdbStore
		err      error
	)
	for {
		if _, err = inp.Read(bSizeBuf); err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		bSize = int(binary.LittleEndian.Uint16(bSizeBuf))
		buf = make([]byte, bSize)
		if _, err = inp.Read(buf); err == io.EOF {
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
		if err = s.rs.Set(iterm.Bucket, iterm.Key, iterm.Value); err != nil && err != bolt.ErrBucketNotFound {
			panic("restore data into backend store happends error: " + err.Error())
		}
		// the backend store is boltdb store and the bucket not exists
		if err == bolt.ErrBucketNotFound {
			// create new bucket
			rs, _ = s.rs.(*store.BoltdbStore)
			if err = rs.CreateBucket(iterm.Bucket); err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// StorageSnapshot for raft
type StorageSnapshot struct {
	diskStore store.RiotStorage
}

// Persist data into disk.
// Notice: every record size can not lager than 131072 byte. mybe that is not good design.
func (s *StorageSnapshot) Persist(sink raft.SnapshotSink) error {
	log.Info("Excute StorageSnapshot.Persist ... ")
	defer sink.Close()
	var (
		c     = s.diskStore.Rec()
		data  []byte
		buf   []byte
		bSize uint16
		err   error
	)
	for {
		iterm := <-c
		if iterm.Err == nil {
			if data, err = json.Marshal(iterm); err != nil {
				return err
			}
			bSize = uint16(len(data))
			buf = make([]byte, bSize+2)
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

// Release snapshot
func (s *StorageSnapshot) Release() {
	log.Info("Excute StorageSnapshot.Release ...")
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
