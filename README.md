# Tries and Radix Trees (blog)

## Overview
This article will present two types of Key-Value data structures; the trie and the radix-tree. First we will provide some context on the importance of these data structures along with high-level descriptions to introduce the concepts.

In addition we provide implementations in Go that should be very similar to the data structures and algorithms in Ethereum. This should provide all the low-level detail required to actually use these concepts in practice.

## Key-Value Stores

Key-Value stores are a simple data storage paradigm that allows us to associate some arbitrary data (the value) with an identifier (the key). They provide efficient methods to;

1. insert a new key-value pair,
2. look up the value for a given key.

Data structures providing this functionality are all over the place. They have become a standard feature in programming languages; such as Python's dictionary or Go's map. It is hard to imagine programming without these tools at ones disposal.

Another common use for key-value stores is NoSQL databases, such as Dynamo and Redis. While these databases generally lack the query power of SQL, they are fast and massively scalable. Because of this they have found countless application in big data and real-time applications.

A final example of application is the Ethereum blockchain app platform. Here key-value stores are used to store state.

## Immutability

Before covering the detail of some key-value implementations we should briefly mention the importance of immutable data.

The idea of immutable data is that once a value is assigned it doesn't ever change. Functional programmers have been evangelical of this style of programming for decades as they have found it makes their programs easier to reason about, debug and parallelise.

Immutable data structures allow a program to update large complex data structures (such as key-value stores) efficiently, while leaving the previous state of the data structure unaffected. Additionally, the new state will reuse the memory from the old state if the data has not changed. This makes such immutable data structures memory efficient if all states need to be retained. 

Diagram BLAH demonstrates how tree based structures allow for this efficient update and sharing of memory. The tree **Si** is modified by changing the node **Ni**. Even if the tree contains many nodes we can get to **Ni** efficiently (provided the tree is balanced). Once **Ni** is replaced with **Ni+1** the ancestors also need to be updated to point to the correct child. All the other nodes in the tree are unaffected and do not need to be copied. If the tree is balanced the number of ancestors will be very small, even if the tree contains many nodes.

This immutability is vital to Ethereum's state storage mechanism as previous states need to be retained so transactions can be verified.

## Trie

Lets start to look at concrete implementations of key-value stores. A trie (as in retrieval, re-**trie**-val) is a type of tree based key-value store. The keys are a sequence of elements from some fixed alphabet, for example strings.

A trie has the property that all descendants of a node have keys with a common prefix. This means insert and look up operations will run in O(k) where k is the length of the key.

Lets walk through an example. Create an empty trie IMAGE. The root of the tree is a null pointer. 

Now insert the pair ("dog", 1). Each letter gets a node, and the leaf node contains the value (1).

Now insert the pair ("cat", 2). There is no shared prefix, so the new data branches after the root.

Now insert the pair ("doge", 3). This key shares a prefix with the existing key "dog", so the new node is added as a child to the existing leaf.

Now insert the pair ("canape", 4). This shares a prefix with "cat", so the new nodes are added as children to the existing branch.

### Implementation

Lets implement this data structure. We will assume the keys are sequences of hexadecimal digits (0-15) and the data is a byte array, as is used in Ethereum.

#### Struct

Lets start by creating a struct for the nodes.

```
type node struct {
	children [16]*node
	data     []byte
}
```

We pre-allocate an array to represent the children of a node. This allows lookup and insert functions to descend to the appropriate child without performing comparisons.

We will need a shallow copy function;

```
func (n *node) copy() *node {
	out := node{
		children: n.children,
		data:     make([]byte, len(n.data)),
	}
	copy(out.data, n.data)
	return &out
}
```

Note that we need to explicit copy the slice, but the array of children gets copied without a call to the copy function.

#### Lookup

Now lets make the lookup function. When performing the lookup function we always have 3 cases;

1. Empty trie so return no data,
2. There is no remaining key to search so return the data at the current node,
3. Otherwise descend to the appropriate child and continue searching.

```
func lookup(n *node, k []byte) []byte {
	if n == nil {
		return nil
	} else if len(k) == 0 {
		return n.data
	} else {
		return lookup(n.children[k[0]], k[1:])
	}
}
```

Pretty simple right?

### Insert

Now lets make the insert function. When performing the insert function we also have 3 cases;

1. Empty trie so add the next node here,
2. There is no remaining key so add the value to the current node,
3. Otherwise descend to the appropriate child and continue.

```
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
```

### Public Interface

Lets hide the implementation behind a public interface.

```
type KVStore interface {
	// This inserts the new (key, value) pair into
	// the trie. k must contain values in [0, 15].
	Insert(k, b []byte) Trie

	// Search the trie for the value associated with
	// the given key. k must contain values in [0, 15].
	// nil is returned if the key is not found.
	Lookup(k []byte) []byte
}

func (n *node) Lookup(k []byte) []byte {
	return lookup(n, k)
}

func (n *node) Insert(k, v []byte) KVStore {
	return insert(n, k, v)
}

func New() *trie {
	return nil
}

type trie node
```

### Example Usage

Let's write a few basic tests. We will check the basic insert and look up functionality. We will also ensure that the structure is immutable.

```
func TestTrie(t *testing.T) {
	t0 := New()
	if t0.Lookup([]byte{0}) != nil {
		t.Errorf("Trie contains data for missing key.")
	}

	t1 := t0.Insert([]byte{0}, []byte("test"))
	if t0.Lookup([]byte{0}) != nil {
		t.Errorf("Insert is not immutable")
	}

	if string(t1.Lookup([]byte{0})) != "test" {
		t.Errorf("Inserted value not found by look up.")
	}

	t2 := t1.Insert([]byte{1}, []byte("another test"))
	if string(t1.Lookup([]byte{0})) != "test" {
		t.Errorf("Insert is not immutable")
	}

	if t1.Lookup([]byte{1}) != nil {
		t.Errorf("Insert is not immutable")
	}

	if string(t2.Lookup([]byte{0})) != "test" {
		t.Errorf("Inserted value not found by look up.")
	}
	if string(t2.Lookup([]byte{1})) != "another test" {
		t.Errorf("Inserted value not found by look up.")
	}

	t3 := t2.Insert([]byte{0, 1}, []byte("a final test"))
	if string(t3.Lookup([]byte{0})) != "test" {
		t.Errorf("Inserted value not found by look up.")
	}
	if string(t3.Lookup([]byte{1})) != "another test" {
		t.Errorf("Inserted value not found by look up.")
	}
	if string(t3.Lookup([]byte{0, 1})) != "a final test" {
		t.Errorf("Inserted value not found by look up.")
	}

}
```

## Radix Tree

The simple trie implementation provided above suffers from at least on glaring issue; memory usage. Every node allocates enough memory to point to 16 children regarless of how many children it has. 

It is possible to trade some time efficiency for space efficiency by replacing the slice of children with a hash map or list of (prefix, node) pairs. Even so a long key with no shared prefix will require each value in the key to be represented by a unique node.

FIGURE

A common optimisation is to compress the trie by merging nodes that have only one child. If we do so the example would look like this.

FIGURE

Now only the branch points need to allocate memory for children, and long sequence of key data can be stored optimally.

### Implementation

A radix tree implementation is signifigantly more complex than a trie. We will base this implementation on that found in Ethereum.

#### Structs

We use several different structs to represent different possible nodes.

They are as follows;
```
type (
	// Just data, no key
	valueNode []byte

	// A compressed prefix shared by multiple values
	shortNode struct {
		key   []byte
		child interface{}
	}

	// A branching node
	fullNode struct {
		value    interface{}
		children [16]interface{}
	}
)

// A shallow copy of a branching node
func (n *fullNode) copy() *fullNode {
	return &fullNode{
		value:    n.value,
		children: n.children,
	}
}
```

You can see this data structure doesn't directly map to the high-level example given above. It still benifits from the compression of shared prefixes and reduction of null children. This diagram demonstrates how we would represent the previous example concretly.

FIGURE

#### Lookup

The look up function has to deal with each of the new node types;

1. The radix tree is empty; return nil
2. The node is a valueNode
	a. The key is empty; return the value
3. The node is a fullNode
	a. The key is empty; return the value at this node
	b. Otherwise call lookup on the appropriate child node
4. The node is a shortNode
	a. There is no shared prefix of the keys; return nil
	b. The shortNode.Key is a prefix of key; call insert on child

```
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
```

#### Insert

The look up now has to deal with more cases;

1. The radix tree is empty
	a. The key is empty; insert a value node
	b. Insert a shortNode with the key, call insert on the child
2. The node is a valueNode
	a. The key is empty; replace the value
	b. Create a branch with the existing data as a child, call insert on the branch
3. The node is a fullNode
	a. The key is empty; insert on the value child
	b. Call insert on the appropriate child node
4. The node is a shortNode
	a. The shortNode.Key is a prefix of key; call insert on child
	b. There is no shared prefix of the keys; create a branch and insert short nodes child and the new (key, value) pair on the branch
	c. Create a branch as above, but insert it under a shortNode containing the shared prefix.

```
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
```

## Benchmarks

## Conclusion

We have demonstrated an implementation of Key-Value stores that is very similar to that featured in Ethereum. For the sake of simplicity we omitted;
1. merkle proofs,
2. encoding of keys,
3. storing nodes on disk for huge Key-Value stores.

These features can be added to the code here with minimal modification.

It is surprisingly easy to implement performant and immutable Key-Value stores. Tires are offer an elegant implementation, but suffer from excessive memory usage. Radix-trees introduce some additional complexity, but the savings in memory usage are significant. 

