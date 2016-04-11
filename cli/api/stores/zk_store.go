package stores

import (
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZkKVIFace interface {
	Children(path string) ([]string, *zk.Stat, error)
	Get(path string) ([]byte, *zk.Stat, error)
	Set(path string, data []byte, version int32) (*zk.Stat, error)
	Delete(path string, version int32) error
	Exists(path string) (bool, *zk.Stat, error)
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
}

type ZkStore struct {
	kv ZkKVIFace
}

func DefaultZkStore() (StoreIFace, error) {
	endpoints := []string{}
	timeout := time.Duration(0)

	client, _, err := zk.Connect(endpoints, timeout)

	if err != nil {
		return nil, err
	}

	return &ZkStore{
		kv: client,
	}, nil
}

func NewZkStore(zn ZkKVIFace) StoreIFace {
	return &ZkStore{
		kv: zn,
	}
}

func (zs *ZkStore) Get(key string) (*KVByte, error) {
	value, _, err := zs.kv.Get(key)

	k := &KVByte{}

	if err != nil || value == nil {
		return nil, err
	}

	k.Key = key
	k.Bytes = value

	return k, nil
}

func (zs *ZkStore) Set(key string, bts []byte) error {
	exists, _, err := zs.kv.Exists(key)
	if err != nil {
		return err
	}

	if !exists {
		zs.createFullPath(splitKey(strings.TrimSuffix(key, "/")))
	}

	_, err = zs.kv.Set(key, bts, -1)

	return err
}

// createFullPath creates the full path for a directory that does not exist
func (zs *ZkStore) createFullPath(path []string) error {
	for i := 1; i <= len(path); i++ {
		p := "/" + strings.Join(path[:i], "/")
		_, err := zs.kv.Create(p, []byte{}, 0, zk.WorldACL(zk.PermAll))

		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

// SplitKey splits the key to extract path information
func splitKey(key string) (path []string) {
	if strings.Contains(key, "/") {
		path = strings.Split(key, "/")
	} else {
		path = []string{key}
	}
	return path
}

func (zs *ZkStore) Delete(key string) error {
	err := zs.kv.Delete(key, -1)

	return err
}

func (zs *ZkStore) List(prefix string) (KVBytes, error) {
	keys, _, err := zs.kv.Children(prefix)

	kvb := make(KVBytes, len(keys))

	if err != nil {
		return kvb, err
	}

	// TODO: Do this in a less-expensive way.
	for i, key := range keys {
		kv, gErr := zs.Get(strings.TrimSuffix(prefix, "/") + key)

		if gErr != nil {
			return nil, gErr
		}

		kvb[i] = &KVByte{
			Key:   key,
			Bytes: kv.Bytes,
		}
	}

	return kvb, err
}

func (zs *ZkStore) Put(key string, bts []byte) error {
	exists, _, err := zs.kv.Exists(key)
	if err != nil {
		return err
	}

	if !exists {
		zs.createFullPath(splitKey(strings.TrimSuffix(key, "/")))
	}

	_, err = zs.kv.Set(key, bts, -1)

	return err
}
