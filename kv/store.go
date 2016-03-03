package kv

import "errors"

var (
	TypeChangeError = errors.New("cannot change existing feature types.")
)

type StoreIFace interface {
	List(prefix string) ([][]byte, error)
	Set(key string, bts []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Put(key string, bts []byte) error
}
