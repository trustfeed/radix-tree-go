# radix-tree-go

A naive implementation of radix-trees, aiming for clarity not performance.

The primary difference between this implementation and the Ethereum implementation is the lack of a node database. The state on the Ethereum blockchain cannot fit into memory so the radix tree needs to provide an efficient method to store on disk, only loading the needed parts of the tree into memory. This complicates the implementation distracting from the radix tree algorithms, so is not featured here.
