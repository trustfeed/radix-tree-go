package pkg

// Make a radix tree
func New() KVStore {
	return &tree{nil}
}

type tree struct {
	root interface{}
}

func (t *tree) Lookup(key []byte) []byte {
	return lookup(t.root, key)
}

func (t *tree) Insert(key, value []byte) KVStore {
	return &tree{insert(t.root, key, value)}
}
