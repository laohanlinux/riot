package store

import (
	"errors"

	"github.com/laohanlinux/go-logger/logger"
)

var ErrFinished = errors.New("all data is sent successfully")
var ErrNotExistBucket = errors.New("the bucket not exists")

const (
	LevelDBStoreBackend = "leveldb"
	BoltDBStoreBackend  = "boltdb"
)

const defaultBucket = "0"

type RiotStorage interface {
	// bucket, key
	Get([]byte, []byte) ([]byte, error)
	// bucket, key, value
	Set([]byte, []byte, []byte) error
	// bucket, key
	Del([]byte, []byte) error
	Rec() <-chan Iterm
}

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
		return NewBoltdbStore(storePath, []byte(defaultBucket))
	default:
		logger.Fatal("unkown the store backend:", storeBackend)
	}
	return nil
}

type Iterm struct {
	Err    error
	Bucket []byte
	Key    []byte
	Value  []byte
}
