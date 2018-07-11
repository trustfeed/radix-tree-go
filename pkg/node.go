package pkg

import (
	"fmt"
	"strings"
)

// This represents a node in the radix tree
type node interface {
}

type (
	// A branching point
	branch struct {
		children [16]node
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
func lookup(node node, key []byte) []byte {
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
			panic("unknown node type")
		}
	}
}

// An iterative implementation of the insert
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
					b.children[n.prefix[plen]] = &leaf{n.prefix[plen+1:], n.data}
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
			plen := prefixLen(key, n.prefix)
			if plen == len(n.prefix) {
				// Keep this node as is and insert to child
				k := n.prefix
				g := f
				f = func(n node) node {
					return g(&compressed{k, n})
				}

				thisNode = n.child
				key = key[plen:]
			} else {
				// Introduce a new branch
				b := branch{}
				if len(n.prefix) > plen+1 {
					b.children[n.prefix[plen]] = &compressed{n.prefix[plen+1:], n.child}
				} else {
					b.children[n.prefix[plen]] = n.child
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
				thisNode = &b
				key = key[plen:]
			}

		default:
			panic("unknown node type")
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

func prettyPrint(origNode node) string {
	switch n := origNode.(type) {
	case nil:
		return "<nil>"
	case *branch:
		strs := make([]string, len(n.children))
		for i, c := range n.children {
			strs[i] = fmt.Sprintf(" %v -> %s ", i, prettyPrint(c))
		}
		return fmt.Sprintf("< %s >", strings.Join(strs, " | "))
	case *leaf:
		return fmt.Sprintf("{ %v -> %s }", n.prefix, n.data)
	case *compressed:
		return fmt.Sprintf("{ %v -> %s }", n.prefix, prettyPrint(n.child))
	default:
		return ""
	}
}
