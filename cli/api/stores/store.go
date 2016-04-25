package stores

import "fmt"

type KVByte struct {
	Key   string
	Bytes []byte
}

func (kv *KVByte) String() string {
	return fmt.Sprintf("%s: %s ", kv.Key, kv.Bytes)
}

type KVBytes []*KVByte

type IFace interface {
	List(prefix string) (KVBytes, error)
	Get(key string) (*KVByte, error)
	Delete(key string) error
	Set(key string, bts []byte) error
	Register(func(kvb KVBytes))
	Watch() error
	Updated(kvs interface{})
	Close()
}
