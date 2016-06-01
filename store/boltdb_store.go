package store

import (
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type BoltdbStore struct {
	*bolt.DB
	c             chan Iterm
	l             *sync.Mutex
	defaultBucket []byte
}

func NewBoltdbStore(dir string, defaultBucket []byte) *BoltdbStore {
	db, err := bolt.Open(dir, 0600, nil)
	if err != nil {
		panic(err)
	}
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

	return &BoltdbStore{
		DB:            db,
		c:             make(chan Iterm),
		l:             &sync.Mutex{},
		defaultBucket: defaultBucket,
	}
}
func (bdbs *BoltdbStore) CreateBucket(bucket []byte) error {
	tx, err := bdbs.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.CreateBucketIfNotExists(bucket)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (bdbs *BoltdbStore) DelBucket(bucket []byte) error {
	tx, err := bdbs.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = tx.DeleteBucket(bucket); err != nil {
		return err
	}
	return tx.Commit()
}

func (bdbs *BoltdbStore) GetBucket(bucket []byte) (bolt.BucketStats, error) {
	var bStats bolt.BucketStats
	tx, err := bdbs.Begin(true)
	if err != nil {
		return bStats, err
	}
	defer tx.Rollback()
	bt := tx.Bucket(bucket)
	if bt != nil {
		return bStats, ErrNotExistBucket
	}
	bStats = bt.Stats()
	return bStats, nil
}

// without transaction
func (bdbs *BoltdbStore) Get(bucket, key []byte) ([]byte, error) {
	var value []byte
	err := bdbs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(key)
		if v != nil {
			value = make([]byte, len(v))
			copy(value, v)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, errors.ErrNotFound
	}
	return value, nil
}

func (bdbs *BoltdbStore) Set(bucket, key, value []byte) error {
	tx, err := bdbs.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	bt := tx.Bucket(bucket)
	if bt == nil {
		return fmt.Errorf("the bucket is not exits")
	}
	if err := bt.Put(key, value); err != nil {
		return err
	}
	return tx.Commit()
}

func (bdbs *BoltdbStore) Del(bucket, key []byte) error {
	tx, err := bdbs.Begin(true)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	bt := tx.Bucket(bucket)
	if err := bt.Delete(key); err != nil {
		return err
	}
	return tx.Commit()
}

func (bdbs *BoltdbStore) Close() error {
	return bdbs.DB.Close()
}

func (bdbs *BoltdbStore) Rec() <-chan Iterm {
	bdbs.l.Lock()
	return bdbs.c
}

func (bdbs *BoltdbStore) streamWorker() {
	bdbs.l.Lock()
	defer bdbs.l.Unlock()
	var iterm Iterm
	bdbs.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bdbs.defaultBucket)
		c := bucket.Cursor()
		iterm.Err = nil
		for k, v := c.First(); k != nil; k, v = c.Next() {
			iterm.Key, iterm.Value = k, v
			bdbs.c <- iterm
		}
		return nil
	})

	iterm.Err = ErrFinished
	iterm.Key, iterm.Value = nil, nil
	bdbs.c <- iterm
}
