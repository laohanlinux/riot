package store

import (
	"errors"
)

var ErrFinished = errors.New("all data is sent successfully")
var ErrNotExistBucket = errors.New("the bucket not exists")

type Iterm struct {
	Err    error
	Bucket []byte
	Key    []byte
	Value  []byte
}
