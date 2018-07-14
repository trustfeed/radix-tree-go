package pkg

import "fmt"

type (
	valueNode []byte

	shortNode struct {
		key   []byte
		child interface{}
	}

	fullNode struct {
		value    interface{}
		children [16]interface{}
	}
)

func (n *fullNode) copy() *fullNode {
	return &fullNode{
		value:    n.value,
		children: n.children,
	}
}

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
			// Keep this node as is and insert to child
			child := insert(n.child, key[plen:], value)
			return &shortNode{n.key, child}

		} else {
			// Introduce a new branch
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
