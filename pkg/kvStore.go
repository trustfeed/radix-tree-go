package pkg

// This is an interface to a Key-Value store
type KVStore interface {
	Insert(key, value []byte) KVStore
	Lookup(key []byte) []byte
}
