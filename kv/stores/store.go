package stores

type StoreIFace interface {
	List(prefix string) ([][]byte, error)
	Get(key string) ([]byte, error)
	Delete(key string) error
	Put(key string, bts []byte) error
}
