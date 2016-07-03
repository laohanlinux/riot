package store

import (
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/laohanlinux/assert"
	"github.com/laohanlinux/go-logger/logger"
)

func TestBoltdbStore(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	benSize := 1 << 10
	logger.Info(benSize)
	dbPath := "tmpBoltdbStoreDir"
	os.RemoveAll(dbPath)
	db := NewBoltdbStore(dbPath)
	assert.NotNil(t, db)
	defer os.RemoveAll(dbPath)
	defer db.Close()
	tBucketName := []byte("Student")

	// 1. bucket test
	assert.Nil(t, db.CreateBucket(tBucketName))
	assert.Nil(t, db.CreateBucket(tBucketName))
	// not exists bucket
	_, err := db.GetBucket([]byte("Not Exists"))
	assert.Equal(t, err.Error(), bolt.ErrBucketNotFound.Error())
	// exists bucket
	_, err = db.GetBucket(tBucketName)
	assert.Nil(t, err)
	// delete bucket
	assert.Nil(t, db.DelBucket(tBucketName))

	// 2. key/value test
	tKey := []byte("zhanShan")
	tValue := []byte(`{"Age":17, "Addr": "Beijin"}`)
	assert.Nil(t, db.CreateBucket(tBucketName))
	value, err := db.Get([]byte("Not Exists"), tKey)
	assert.Nil(t, value)
	assert.Equal(t, err.Error(), bolt.ErrBucketNotFound.Error())

	err = db.Set([]byte("Not Exists"), tKey, tValue)
	assert.Equal(t, err, bolt.ErrBucketNotFound)
	assert.Nil(t, db.Set(tBucketName, tKey, tValue))
	value, err = db.Get(tBucketName, tKey)
	assert.Nil(t, err)
	assert.Equal(t, value, tValue)
	assert.Nil(t, db.Del(tBucketName, tKey))

	// 3. test scan all data
	assert.Nil(t, db.DelBucket(tBucketName))
	assert.Nil(t, db.CreateBucket(tBucketName))
	for k := 0; k < benSize; k++ {
		v := k + 1
		assert.Nil(t, db.Set(tBucketName, []byte(strconv.Itoa(k)), []byte(strconv.Itoa(v))))
	}

	c := db.Rec()

	for iterm := range c {
		if iterm.Err == ErrFinished {
			break
		}
		key, err := strconv.Atoi(string(iterm.Key))
		assert.Nil(t, err)
		value, err := strconv.Atoi(string(iterm.Value))
		assert.Nil(t, err)
		assert.Equal(t, key+1, value)
	}
}
