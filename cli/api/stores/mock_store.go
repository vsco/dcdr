package stores

import "github.com/vsco/dcdr/models"

type MockStore struct {
	Item  *KVByte
	Items KVBytes
	Err   error
}

func NewMockStore(ft *models.Feature, err error) (ms *MockStore) {
	bts, _ := ft.ToJSON()

	ms = &MockStore{
		Err: err,
	}

	if ft != nil {
		kvb := KVBytes{
			&KVByte{
				Key:   ft.Key,
				Bytes: bts,
			},
		}

		ms.Item = kvb[0]
		ms.Items = kvb
	}

	return
}

func (ms *MockStore) List(prefix string) (KVBytes, error) {
	return ms.Items, ms.Err
}

func (ms *MockStore) Get(key string) (*KVByte, error) {
	return ms.Item, ms.Err
}

func (ms *MockStore) Set(key string, bts []byte) error {
	return ms.Err
}

func (ms *MockStore) Delete(key string) error {
	return ms.Err
}

func (ms *MockStore) Put(key string, bts []byte) error {
	return ms.Err
}

type MockRepo struct {
	error   error
	sha     string
	exists  bool
	enabled bool
}

func (mr *MockRepo) Clone() error {
	return mr.error
}

func (mr *MockRepo) Commit(bts []byte, msg string) error {
	return mr.error
}

func (mr *MockRepo) Create() error {
	return mr.error
}

func (mr *MockRepo) Exists() bool {
	return mr.exists
}

func (mr *MockRepo) Enabled() bool {
	return mr.enabled
}

func (mr *MockRepo) Push() error {
	return mr.error
}

func (mr *MockRepo) Pull() error {
	return mr.error
}

func (mr *MockRepo) CurrentSHA() (string, error) {
	return mr.sha, mr.error
}

func (mr *MockRepo) Init() {
}
