package fsm

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/laohanlinux/riot/rpc/pb"
	"github.com/laohanlinux/riot/share"
	"github.com/laohanlinux/riot/store"

	"github.com/hashicorp/raft"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var ErrNotFound = fmt.Errorf("the key's value is nil.")
var ErrInvalidCmd = fmt.Errorf("The command is invalid")

func NewStorageFSM(rs RiotStorage) *StorageFSM {
	return &StorageFSM{
		l:  &sync.Mutex{},
		rs: rs,
	}
}

// StorageFSM is an implememtation of the FSM interfacec, and just
// storage the key/value logs sequentially
type StorageFSM struct {
	l  *sync.Mutex
	rs RiotStorage
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
	case "SET":
		err = s.rs.Set([]byte(req.Key), req.Value)
	case "DEL":
		err = s.rs.Del([]byte(req.Key))
	case "SHARE":
		addr := string(req.Value)
		idx := strings.Index(addr, ":")
		share.ShCache.LRPC.Addr = addr[:idx]
		share.ShCache.LRPC.Port = addr[idx+1:]
		//logger.Debug("update share cache meory:", string(req.Value), share.ShCache.LRPC.Addr, share.ShCache.LRPC.Port)
	default:
		err = ErrInvalidCmd
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// Get .
func (s *StorageFSM) Get(key string) ([]byte, error) {
	s.l.Lock()
	defer s.l.Unlock()
	value, err := s.rs.Get([]byte(key))
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
	//logger.Info("Excute StorageFSM.Snapshot ...")
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
			panic(err)
		}
		// decoding
		if err = json.Unmarshal(buf, &iterm); err != nil {
			panic(err)
		}
		if err = s.rs.Set(iterm.Key, iterm.Value); err != nil {
			panic(err)
		}
	}

	return nil
}

// StorageSnapshot .
type StorageSnapshot struct {
	diskCache map[string][]byte
	diskStore RiotStorage
}

// Persist ...
func (s *StorageSnapshot) Persist(sink raft.SnapshotSink) error {
	//logger.Info("Excute StorageSnapshot.Persist ... ")
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
