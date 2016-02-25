package command

import (
	"github.com/laohanlinux/riot/fsm"
)

const (
	cmdGet = "get"
	cmdSet = "set"
	cmdDel = "del"
)

type Comand struct {
	op      string
	raftFSM *fsm.StorageFSM
}

func (cm Comand) doGet(key []byte) []byte {
	return nil
}

func (cm Comand) doSet(key, value []byte) error {
	return nil
}

func (cm Comand) doDel(key []byte) error {
	return nil
}
