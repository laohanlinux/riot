package store

import (
	"sync"

	"github.com/boltdb/bolt"
	log "github.com/laohanlinux/utils/gokitlog"
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
	var (
		tx  *bolt.Tx
		err error
	)
	if tx, err = bdbs.Begin(true); err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.CreateBucketIfNotExists(bucket); err != nil {
		return err
	}

	bdbs.l.Lock()

	if _, ok := bdbs.buckets[string(bucket)]; !ok {
		bdbs.buckets[string(bucket)] = true
	}
	bdbs.l.Unlock()

	return tx.Commit()
}

func (bdbs *BoltdbStore) DelBucket(bucket []byte) (err error) {
	var (
		tx *bolt.Tx
		ok bool
	)
	if tx, err = bdbs.Begin(true); err != nil {
		return err
	}
	defer tx.Rollback()
	if err = tx.DeleteBucket(bucket); err != nil {
		return err
	}

	bdbs.l.Lock()
	if _, ok = bdbs.buckets[string(bucket)]; ok {
		delete(bdbs.buckets, string(bucket))
	}
	bdbs.l.Unlock()

	return tx.Commit()
}

func (bdbs *BoltdbStore) GetBucket(bucket []byte) (bStats bolt.BucketStats, err error) {
	var (
		tx *bolt.Tx
		bt *bolt.Bucket
	)
	if tx, err = bdbs.Begin(true); err != nil {
		return bStats, err
	}
	defer tx.Rollback()
	if bt = tx.Bucket(bucket); bt == nil {
		return bStats, bolt.ErrBucketNotFound
	}
	bStats = bt.Stats()
	return
}

// Get implements the RiotStorage interface
func (bdbs *BoltdbStore) Get(bucket, key []byte) (value []byte, err error) {

	err = bdbs.View(func(tx *bolt.Tx) error {
		var (
			bt *bolt.Bucket
			v  []byte
		)
		if bt = tx.Bucket(bucket); bt == nil {
			return bolt.ErrBucketNotFound
		}
		if v = bt.Get(key); v != nil {
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
	var (
		tx  *bolt.Tx
		bt  *bolt.Bucket
		err error
	)
	if tx, err = bdbs.Begin(true); err != nil {
		return err
	}
	defer tx.Rollback()
	if bt = tx.Bucket(bucket); bt == nil {
		return bolt.ErrBucketNotFound
	}
	if err = bt.Put(key, value); err != nil {
		return err
	}
	return tx.Commit()
}

// Del implements the RiotHandler interface
func (bdbs *BoltdbStore) Del(bucket, key []byte) error {
	var (
		tx  *bolt.Tx
		bt  *bolt.Bucket
		err error
	)
	if tx, err = bdbs.Begin(true); err != nil {
		return nil
	}
	defer tx.Rollback()

	if bt = tx.Bucket(bucket); bt == nil {
		return bolt.ErrBucketNotFound
	}

	if err = bt.Delete(key); err != nil {
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
			var (
				bt *bolt.Bucket
				c  *bolt.Cursor
			)
			if bt = tx.Bucket([]byte(bucketName)); bt == nil {
				log.Error("bucket", nil)
				continue
			}
			c = bt.Cursor()
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
