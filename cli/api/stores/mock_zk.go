package stores

import "github.com/samuel/go-zookeeper/zk"

type MockZk struct {
	Item  *KVByte
	Items KVBytes
	Err   error
}

func NewMockZk(key string, kvb KVBytes, err error) (mc *MockZk) {
	mc = &MockZk{
		Err: err,
	}

	if len(kvb) != 0 {
		mc.Item = kvb[0]
		mc.Items = kvb
	}

	return
}

func (mc *MockZk) Children(path string) ([]string, *zk.Stat, error) {
	items := []string{mc.Item.Key}

	return items, nil, mc.Err
}

func (mc *MockZk) Get(path string) ([]byte, *zk.Stat, error) {
	return mc.Item.Bytes, nil, mc.Err
}

func (mc *MockZk) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	return nil, mc.Err
}

func (mc *MockZk) Delete(path string, version int32) error {
	return mc.Err
}

func (mc *MockZk) Exists(path string) (bool, *zk.Stat, error) {
	return true, nil, mc.Err
}

func (mc *MockZk) Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return path, mc.Err
}
