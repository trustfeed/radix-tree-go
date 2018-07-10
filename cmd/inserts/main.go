package main

import (
	"fmt"

	radix "github.com/trustfeed/radix-tree-go/pkg"
)

func main() {
	t := radix.Make()
	t1 := t.Insert([]byte{0}, []byte{10})
	t2 := t1.Insert([]byte{0}, []byte{11})
	t3 := t1.Insert([]byte{0, 1}, []byte{15})

	fmt.Println(t.Lookup([]byte{0}))
	fmt.Println(t1.Lookup([]byte{0}))
	fmt.Println(t2.Lookup([]byte{0}))
	fmt.Println(t3.Lookup([]byte{0}))
	fmt.Println(t3.Lookup([]byte{0, 1}))
}
