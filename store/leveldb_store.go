package store

import (
	"sync"

	"github.com/laohanlinux/go-logger/logger"
	"github.com/syndtr/goleveldb/leveldb"
)

type LeveldbStorage struct {
	*leveldb.DB
	c chan Iterm
	l *sync.Mutex
}

// NewLeveldbStorage returns a dbstore
func NewLeveldbStorage(dir string) *LeveldbStorage {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		panic(err)
	}
	return &LeveldbStorage{
		c:  make(chan Iterm),
		l:  &sync.Mutex{},
		DB: db,
	}
}

func (edbs *LeveldbStorage) Get(_, key []byte) ([]byte, error) {
	logger.Info("Get a value by ", string(key))
	return edbs.DB.Get(key, nil)
}

func (edbs *LeveldbStorage) Set(_, key, value []byte) error {
	logger.Info("Set a key/value:", string(key), string(value))
	return edbs.DB.Put(key, value, nil)
}

func (edbs *LeveldbStorage) Del(_, key []byte) error {
	return edbs.DB.Delete(key, nil)
}

func (edbs *LeveldbStorage) Close() error {
	return edbs.DB.Close()
}

func (edbs *LeveldbStorage) Rec() <-chan Iterm {
	edbs.l.Lock()
	go edbs.streamWorker()
	return edbs.c
}

// TODO:
//
// don't alloc new memory for every time,
// because the iter.Key and iter.Value use same memory space every iter.
func (edbs *LeveldbStorage) streamWorker() {
	defer edbs.l.Unlock()
	iter := edbs.NewIterator(nil, nil)
	var iterm Iterm
	for iter.Next() {
		iterm.Err, iterm.Bucket = nil, nil
		iterm.Key = make([]byte, len(iter.Key()))
		iterm.Value = make([]byte, len(iter.Value()))
		copy(iterm.Key, iter.Key())
		copy(iterm.Value, iter.Value())
		logger.Debug(string(iterm.Key), string(iterm.Value))
		edbs.c <- iterm
	}
	iterm.Key, iterm.Value = nil, nil
	iterm.Err = ErrFinished
	edbs.c <- iterm
}
