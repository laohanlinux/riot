package store

import (
	"errors"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/boltdb/bolt"
	"fmt"
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
	go edbs.streamWorker()
	return edbs.c
}

func (edbs *leveldbStorage) streamWorker() {
	edbs.l.Lock()
	defer edbs.l.Unlock()
	iter := edbs.NewIterator(nil, nil)
	var iterm Iterm
	for iter.Next() {
		iterm.Err = nil
		iterm.Key, iterm.Value = iter.Key(), iter.Value()
		edbs.c <- iterm
	}
	iterm.Key, iterm.Value = nil, nil
	iterm.Err = ErrFinished
	edbs.c <- iterm
}


// boltdb store
type boltdbStore  struct{
	*bolt.DB
	c chan Iterm
	l *sync.Mutex
	defaultBucket []byte
}

func NewBoltdbStore(dir string, defaultBucket []byte) *boltdbStore{
	db, err := bolt.Open(dir, 0600, nil)
	if err != nil {
		panic(err)
	}
	// create a new bucket
	tx, err := db.Begin(true)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	_, err = tx.CreateBucketIfNotExists(defaultBucket)
	if err != nil {
		panic(err)
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return &boltdbStore{
		DB: db,
		c : make(chan Iterm),
		l : &sync.Mutex{},
		defaultBucket: defaultBucket,
	}
}
// without transaction
func (bdbs * boltdbStore) Get(key []byte)([]byte, error) {
	var value []byte
	err := bdbs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bdbs.defaultBucket)
		v := b.Get(key)
		if v != nil {
			copy(value, v)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (bdbs *boltdbStore) Set(key, value []byte) error {
	tx, err := bdbs.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	fmt.Println("key:", string(key), "value:", value, "bucket:", string(bdbs.defaultBucket))
	bucket := tx.Bucket(bdbs.defaultBucket)
	if bucket == nil {
		return fmt.Errorf("the bucket is nil.")
	}
	if err := bucket.Put(key, value); err != nil {
		return err
	}
	return tx.Commit()
}

func (bdbs * boltdbStore) Del(key []byte) error {
	tx, err := bdbs.Begin(true)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bdbs.defaultBucket)
	if err := bucket.Delete(key); err != nil {
		return err
	}
	return tx.Commit()
}

func (bdbs *boltdbStore) Close() error {
	return bdbs.DB.Close()
}

func (bdbs *boltdbStore) Rec() <-chan Iterm {
	bdbs.l.Lock()
	return bdbs.c
}

func (bdbs *boltdbStore) streamWorker() {
	bdbs.l.Lock()
	defer  bdbs.l.Unlock()
	var iterm Iterm
	bdbs.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket(bdbs.defaultBucket)
		c := bucket.Cursor()
		iterm.Err = nil
		for k, v := c.First(); k != nil; k, v = c.Next(){
			iterm.Key, iterm.Value = k, v
			bdbs.c <- iterm
		}
		return nil
	})

	iterm.Err = ErrFinished
	iterm.Key, iterm.Value = nil, nil
	bdbs.c <- iterm
}

