package fsm

import (
	"github.com/laohanlinux/riot/store"

	"github.com/laohanlinux/go-logger/logger"
)

const (
	LevelDBStoreBackend = "leveldb"
	BoltDBStoreBackend = "boltdb"
)

const defaultBucket  = "0"
type RiotStorage interface {
	// bucket, key
	Get([]byte, []byte) ([]byte, error)
	// bucket, key, value
	Set([]byte, []byte, []byte) error
	// bucket, key
	Del([]byte, []byte) error
	Rec() <-chan store.Iterm
}

type RiotStorageFactory struct {}

var rsf *RiotStorageFactory

// NewRiotStoreFactory is not a thread safely function.
func NewRiotStoreFactory(storeBackend, storePath string) RiotStorage {
	if rsf == nil {
		rsf = &RiotStorageFactory{}
	}
	switch storeBackend {
	case LevelDBStoreBackend:
		return store.NewLeveldbStorage(storePath)
	case BoltDBStoreBackend:
		return store.NewBoltdbStore(storePath, []byte(defaultBucket))
	default:
		logger.Fatal("unkown the store backend:", storeBackend)
	}
	return nil
}
