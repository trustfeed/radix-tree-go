package pkg

import "reflect"

// This represents a node in the radix tree
type node interface {
	nothing()
}

func (*branch) nothing()     {}
func (*leaf) nothing()       {}
func (*compressed) nothing() {}

type (
	// A branching point
	branch struct {
		children [17]node
		data     []byte
	}

	// No children below here
	leaf struct {
		prefix []byte
		data   []byte
	}

	// A sequence of data with no branches
	compressed struct {
		prefix []byte
		child  node
	}
)

// An iterative implementation of lookup function.
func lookup(originalNode node, key []byte) []byte {
	node := originalNode
	for {
		switch n := node.(type) {
		case nil:
			return nil

		case *branch:
			if len(key) == 0 {
				return n.data
			}
			node = n.children[key[0]]
			key = key[1:]

		case *leaf:
			if byteArrayEqual(key, n.prefix) {
				return n.data
			}
			return nil

		case *compressed:
			plen := len(n.prefix)
			if plen > len(key) || !byteArrayEqual(key[:plen], n.prefix) {
				return nil
			}
			node = n.child
			key = key[plen:]

		default:
			panic(reflect.TypeOf(n))
		}
	}
}

func insert(originalNode node, key, value []byte) node {
	f := func(node node) node {
		return node
	}

	thisNode := originalNode
	for {
		switch n := thisNode.(type) {
		case nil:
			return f(&leaf{key, value})

		case *branch:
			if len(key) == 0 {
				b := n.copy()
				b.data = value
				return f(b)
			}

			b := n.copy()
			k := key[0]
			g := f
			f = func(n node) node {
				b.children[k] = n
				return g(b)
			}
			thisNode = n.children[key[0]]
			key = key[1:]

		case *leaf:
			if byteArrayEqual(key, n.prefix) {
				return f(&leaf{key, value})
			} else {
				plen := prefixLen(key, n.prefix)

				b := branch{}
				if plen == len(n.prefix) {
					b.data = n.data
				} else {
					b.children[n.prefix[plen]] = &leaf{n.prefix[plen:], n.data}
				}

				g := f
				k := key[:plen]
				f = func(n node) node {
					if plen == 0 {
						return g(n)
					} else {
						return g(&compressed{k, n})
					}
				}

				key = key[plen:]
				thisNode = &b
			}

		case *compressed:
			return f(nil)

		default:
			panic("what")
		}

	}
}

func (b *branch) copy() *branch {
	out := branch{
		children: b.children,
		data:     make([]byte, len(b.data)),
	}
	copy(out.data, b.data)
	return &out
}
