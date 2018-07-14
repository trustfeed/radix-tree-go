package pkg

type KVStore interface {
	Insert(key, value []byte) KVStore
	Lookup(key []byte) []byte
}
