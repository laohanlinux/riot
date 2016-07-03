package store

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/laohanlinux/assert"
	"github.com/laohanlinux/go-logger/logger"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func TestLeveldbStore(t *testing.T) {
	//	runtime.GOMAXPROCS(runtime.NumCPU())
	benSize := 1 << 10
	logger.Info(benSize)
	dbPath := "tmpLevelDBStoreDir"
	os.RemoveAll(dbPath)
	db := NewLeveldbStorage(dbPath)
	assert.NotNil(t, db)
	defer os.RemoveAll(dbPath)
	defer db.Close()

	// 1. key/value test
	tKey := []byte("zhanShan")
	tValue := []byte(`{"Age":17, "Addr": "Beijin"}`)
	value, err := db.Get(nil, tKey)
	assert.Nil(t, value)
	assert.Equal(t, err.Error(), errors.ErrNotFound.Error())

	err = db.Set(nil, tKey, tValue)
	assert.Nil(t, err)
	value, err = db.Get(nil, tKey)
	assert.Nil(t, err)
	assert.Equal(t, value, tValue)
	assert.Nil(t, db.Del(nil, tKey))
	value, err = db.Get(nil, tKey)
	assert.Equal(t, err, errors.ErrNotFound)
	assert.Equal(t, len(value), 0)

	// 3. test scan all data
	for k := 0; k < benSize; k++ {
		v := k + 1
		assert.Nil(t, db.Set(nil, []byte(strconv.Itoa(k)), []byte(strconv.Itoa(v))))
	}

	c := db.Rec()

	for iterm := range c {
		if iterm.Err == ErrFinished {
			break
		}
		fmt.Printf("%p\t%s\t%s\t%p\t%p\n", &iterm, iterm.Key, iterm.Value, iterm.Key, iterm.Value)
		logger.Info(string(iterm.Key), string(iterm.Value))
		key, err := strconv.Atoi(string(iterm.Key))
		assert.Nil(t, err)
		value, err := strconv.Atoi(string(iterm.Value))
		assert.Nil(t, err)
		assert.Equal(t, key+1, value)
	}
}
