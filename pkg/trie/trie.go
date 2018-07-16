package trie

import radix "github.com/trustfeed/radix-tree-go/pkg"

// A node in a trie.
type node struct {
	children [16]*node
	data     []byte
}

// A shallow copy of a trie node.
func (n *node) copy() *node {
	out := node{
		children: n.children,
		data:     make([]byte, len(n.data)),
	}
	copy(out.data, n.data)
	return &out
}

func lookup(n *node, k []byte) []byte {
	if n == nil {
		return nil
	} else if len(k) == 0 {
		return n.data
	} else {
		return lookup(n.children[k[0]], k[1:])
	}
}

func insert(n *node, k []byte, v []byte) *node {
	if n == nil {
		out := node{}
		return insert(&out, k, v)
	} else if len(k) == 0 {
		out := n.copy()
		out.data = v
		return out
	} else {
		out := n.copy()
		out.children[k[0]] = insert(n.children[k[0]], k[1:], v)
		return out
	}
}

func (n *node) Lookup(k []byte) []byte {
	return lookup(n, k)
}

func (n *node) Insert(k, v []byte) radix.KVStore {
	return insert(n, k, v)
}

func New() *node {
	return nil
}
