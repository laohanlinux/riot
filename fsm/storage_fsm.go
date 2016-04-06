package fsm

import (
	"github.com/laohanlinux/riot/store"
)

type RiotStorage interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte) error
	Del([]byte) error
	Rec() <-chan store.Iterm
}
