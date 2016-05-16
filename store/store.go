package store

import (
	"errors"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var ErrFinished = errors.New("all data is sent successfully")

type leveldbStorage struct {
	*leveldb.DB
	c chan Iterm
	l *sync.Mutex
}

type Iterm struct {
	Err   error
	Key   []byte
	Value []byte
}

// NewLeveldbStorage returns a dbstore
func NewLeveldbStorage(dir string) *leveldbStorage {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		panic(err)
	}
	return &leveldbStorage{
		c:  make(chan Iterm),
		l:  &sync.Mutex{},
		DB: db,
	}
}

func (edbs *leveldbStorage) Get(key []byte) ([]byte, error) {
	return edbs.DB.Get(key, nil)
}

func (edbs *leveldbStorage) Set(key, value []byte) error {
	return edbs.DB.Put(key, value, nil)
}

func (edbs *leveldbStorage) Del(key []byte) error {
	return edbs.DB.Delete(key, nil)
}

func (edbs *leveldbStorage) Close() error {
	return edbs.DB.Close()
}

func (edbs *leveldbStorage) Rec() <-chan Iterm {
	edbs.l.Lock()
	return edbs.c
}

func (edbs *leveldbStorage) streamWorker() {
	defer edbs.l.Unlock()
	iter := edbs.NewIterator(nil, nil)
	var iterm Iterm
	for iter.Next() {
		iterm.Err = nil
		iterm.Key, iterm.Value = iter.Key(), iter.Value()
		edbs.c <- iterm
	}
	iterm.Err = ErrFinished
	edbs.c <- iterm
}


// boltdb store
type boltdbStore  struct{

}