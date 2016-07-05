package store

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type BoltdbStore struct {
	*bolt.DB
	c       chan Iterm
	l       *sync.RWMutex
	buckets map[string]bool
}

func NewBoltdbStore(dir string) *BoltdbStore {
	db, err := bolt.Open(dir, 0600, nil)
	if err != nil {
		panic(err)
	}
	return &BoltdbStore{
		DB:      db,
		c:       make(chan Iterm),
		l:       &sync.RWMutex{},
		buckets: make(map[string]bool),
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

	bdbs.l.Lock()
	if _, ok := bdbs.buckets[string(bucket)]; !ok {
		bdbs.buckets[string(bucket)] = true
	}
	bdbs.l.Unlock()

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

	bdbs.l.Lock()
	if _, ok := bdbs.buckets[string(bucket)]; ok {
		delete(bdbs.buckets, string(bucket))
	}
	bdbs.l.Unlock()

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
	if bt == nil {
		return bStats, bolt.ErrBucketNotFound
	}
	bStats = bt.Stats()
	return bStats, nil
}

// Get implements the RiotStorage interface
func (bdbs *BoltdbStore) Get(bucket, key []byte) ([]byte, error) {
	var value []byte
	err := bdbs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return bolt.ErrBucketNotFound
		}
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

// Set implements the RiotStorage interface
func (bdbs *BoltdbStore) Set(bucket, key, value []byte) error {
	tx, err := bdbs.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	bt := tx.Bucket(bucket)
	if bt == nil {
		return bolt.ErrBucketNotFound
	}
	if err := bt.Put(key, value); err != nil {
		return err
	}
	return tx.Commit()
}

// Del implements the RiotHandler interface
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

// Rec implements the RiotStorage interface
func (bdbs *BoltdbStore) Rec() <-chan Iterm {
	bdbs.l.Lock()
	go bdbs.streamWorker()
	return bdbs.c
}

// Close recycle the db source
func (bdbs *BoltdbStore) Close() error {
	return bdbs.DB.Close()
}

// get all data spends too much time ...
func (bdbs *BoltdbStore) streamWorker() {
	defer bdbs.l.Unlock()
	var iterm Iterm
	bdbs.View(func(tx *bolt.Tx) error {
		// scan all bucket
		for bucketName, _ := range bdbs.buckets {
			bucket := tx.Bucket([]byte(bucketName))
			c := bucket.Cursor()
			iterm.Err = nil
			for k, v := c.First(); k != nil; k, v = c.Next() {
				iterm.Bucket, iterm.Key, iterm.Value = []byte(bucketName), k, v
				bdbs.c <- iterm
			}
		}
		return nil
	})

	iterm.Err = ErrFinished
	iterm.Key, iterm.Value = nil, nil
	bdbs.c <- iterm
}
