package fsm

import "github.com/syndtr/goleveldb/leveldb"

type riotStorage interface {
	open(string)
	get([]byte) ([]byte, error)
	set([]byte, []byte) error
	del([]byte) error
	close() error
}

type leveldbStorage struct {
	*leveldb.DB
}

func (edbs *leveldbStorage) open(dir string) error {
	db, err := leveldb.OpenFile("path/to/db", nil)
	edbs.DB = db
	return err
}

func (edbs *leveldbStorage) get(key []byte) ([]byte, error) {
	return edbs.DB.Get(key, nil)
}

func (edbs *leveldbStorage) set(key, value []byte) error {
	return edbs.DB.Put(key, value, nil)
}

func (edbs *leveldbStorage) del(key []byte) error {
	return edbs.DB.Delete(key, nil)
}

func (edbs *leveldbStorage) close() error {
	return edbs.DB.Close()
}
