package pkg

type Tree interface {
	Lookup(key []byte) []byte
	Insert(key, value []byte) Tree
}

func Make() Tree {
	return &tree{}
}

type tree struct {
	root node
}

func (t *tree) Lookup(key []byte) []byte {
	return lookup(t.root, key)
}

func (t *tree) Insert(key, value []byte) Tree {
	return &tree{insert(t.root, key, value)}
}
