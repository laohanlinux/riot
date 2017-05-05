package store

import (
	"errors"
)

// ErrFinished indicates all data restore
var ErrFinished = errors.New("all data is sent successfully")

const (
	// LevelDBStoreBackend of leveldb store
	LevelDBStoreBackend = "leveldb"
	// BoltDBStoreBackend of boltdb store
	BoltDBStoreBackend = "boltdb"
)

// RiotStorage is a store interface
type RiotStorage interface {
	// bucket, key
	Get([]byte, []byte) ([]byte, error)
	// bucket, key, value
	Set([]byte, []byte, []byte) error
	// bucket, key
	Del([]byte, []byte) error
	Rec() <-chan Iterm
}

// RiotStorageFactory is a store Factory.
type RiotStorageFactory struct{}

var rsf *RiotStorageFactory

// NewRiotStoreFactory is not a thread safely function.
func NewRiotStoreFactory(storeBackend, storePath string) RiotStorage {
	if rsf == nil {
		rsf = &RiotStorageFactory{}
	}
	switch storeBackend {
	case LevelDBStoreBackend:
		return NewLeveldbStorage(storePath)
	case BoltDBStoreBackend:
		return NewBoltdbStore(storePath)
	default:
		panic("unkown the store backend:" + storeBackend)
	}
	return nil
}

type Iterm struct {
	Err    error
	Bucket []byte
	Key    []byte
	Value  []byte
}
