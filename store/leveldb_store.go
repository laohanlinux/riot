package store

import (
	"sync"

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
	return edbs.DB.Get(key, nil)
}

func (edbs *LeveldbStorage) Set(_, key, value []byte) error {
	return edbs.DB.Put(key, value, nil)
}

func (edbs *LeveldbStorage) Del(_, key []byte) error {
	return edbs.DB.Delete(key, nil)
}

func (edbs *LeveldbStorage) Close() error {
	return edbs.DB.Close()
}

func (edbs *LeveldbStorage) Rec() <-chan Iterm {
	go edbs.streamWorker()
	return edbs.c
}

func (edbs *LeveldbStorage) streamWorker() {
	edbs.l.Lock()
	defer edbs.l.Unlock()
	iter := edbs.NewIterator(nil, nil)
	var iterm Iterm
	for iter.Next() {
		iterm.Err, iterm.Bucket = nil, nil
		iterm.Key, iterm.Value = iter.Key(), iter.Value()
		edbs.c <- iterm
	}
	iterm.Key, iterm.Value = nil, nil
	iterm.Err = ErrFinished
	edbs.c <- iterm
}
