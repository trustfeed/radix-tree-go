# radix-tree-go

A naive implementation of radix-trees, aiming for clarity not performance.

The primary difference between this implementation and the Ethereum implementation is the lack of a node database. The state on the Ethereum blockchain cannot fit into memory so the radix tree needs to provide an efficient method to store on disk, only loading the needed parts of the tree into memory. This complicates the implementation distracting from the radix tree algorithms, so is not featured here.

# Data Structures

A radix tree is made of nodes. Each node can be either;

1. An empty tree (nil),
2. A branch with 16 children nodes and optional data,
3. A leaf node with the remaining key and a value, or
4. A compressed prefix with prefix data and a single child node.

# Example

Lets walk through an example;

1. Begin with an empty tree (S<sub>0</sub>).
2. Insert the (key, value) pair ("dog", 1), resulting in a single leaf (S<sub>1</sub>).
3. Insert the (key, value) pair ("cat", 2), resulting in a branch with two leaves as children (S<sub>2</sub>).
4. Insert the (key, value) pair ("doge", 3), resulting in a branch with a leaf and compressed prefix as children (S<sub>3</sub>).

![Example radix trees](https://raw.githubusercontent.com/trustfeed/radix-tree-go/)
