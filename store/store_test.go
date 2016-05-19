package store

import (
	"testing"
	"os"

	"github.com/laohanlinux/assert"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/laohanlinux/go-logger/logger"
)

func TestBoltdbStore(t *testing.T) {
	logger.Info("hello Word")
	boltdbDir := "boltdbDir"
	bucketName := []byte("0")
	// clear dirty files
	os.RemoveAll(boltdbDir)
	defer os.RemoveAll(boltdbDir)
	db := NewBoltdbStore(boltdbDir, bucketName)
	assert.NotNil(t, db)

	testkey, testValue := []byte("Hello"), []byte("Word")
	resValue, err := db.Get(testkey)
	assert.Equal(t, err.Error(), errors.ErrNotFound.Error())
	assert.Nil(t, resValue)

	// set the value
	assert.Nil(t, db.Set(testkey, testValue))

	// get the value
	resValue, err = db.Get(testkey)
	assert.Nil(t, err)
	assert.Equal(t, testValue, resValue)

	// delete the value
	assert.Nil(t, db.Del(testkey))

	assert.Nil(t, db.Close())
}
