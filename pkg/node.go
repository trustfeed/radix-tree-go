package pkg

import "fmt"

type (
	// The value in the radix-tree
	valueNode []byte

	// A compressed key
	shortNode struct {
		key   []byte
		child interface{}
	}

	// Branching node
	fullNode struct {
		value    interface{}
		children [16]interface{}
	}
)

// Shallow copy a branch
func (n *fullNode) copy() *fullNode {
	return &fullNode{
		value:    n.value,
		children: n.children,
	}
}

// Search for a key
func lookup(node interface{}, key []byte) []byte {
	switch n := node.(type) {
	case nil:
		return nil

	case *valueNode:
		if len(key) == 0 {
			return *n
		} else {
			return nil
		}

	case *shortNode:
		plen := len(n.key)
		if plen > len(key) || !byteArrayEqual(key[:plen], n.key) {
			return nil
		} else {
			return lookup(n.child, key[plen:])
		}

	case *fullNode:
		if len(key) == 0 {
			return lookup(n.value, key)
		} else {
			return lookup(n.children[key[0]], key[1:])
		}

	default:
		fmt.Println(n)
		panic(fmt.Sprintf("unknown node type %v", n))
	}
}

// Insert Key-Value pair
func insert(node interface{}, key, value []byte) interface{} {
	switch n := node.(type) {
	case nil:
		if len(key) == 0 {
			out := make(valueNode, len(value))
			copy(out, value)
			return &out
		} else {
			return &shortNode{
				key:   key,
				child: insert(nil, nil, value),
			}
		}

	case *valueNode:
		if len(key) == 0 {
			out := make(valueNode, len(value))
			copy(out, value)
			return &out
		} else {
			b := insert(&fullNode{value: n}, key, value)
			return b
		}

	case *shortNode:
		plen := prefixLen(key, n.key)
		if plen == len(n.key) {
			child := insert(n.child, key[plen:], value)
			return &shortNode{n.key, child}

		} else {
			b := fullNode{}
			if len(n.key) > plen+1 {
				b.children[n.key[plen]] = &shortNode{n.key[plen+1:], n.child}
			} else {
				b.children[n.key[plen]] = n.child
			}

			child := insert(&b, key[plen:], value)

			if plen == 0 {
				return child
			} else {
				return &shortNode{key[:plen], child}
			}
		}

	case *fullNode:
		b := n.copy()

		if len(key) == 0 {
			b.value = insert(b.value, nil, value)
			return b
		} else {
			k := key[0]
			newChild := insert(n.children[k], key[1:], value)
			b.children[k] = newChild
			return b
		}

	default:
		panic(fmt.Sprintf("unknown node type %v", n))
	}
}
