package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	radix "github.com/trustfeed/radix-tree-go/pkg"
	trie "github.com/trustfeed/radix-tree-go/pkg/trie"
)

type input struct {
	key   []byte
	value []byte
}

func readLine(file *bufio.Reader) (input, bool) {
	out := input{}

	{
		l, err := file.ReadString('\n')
		if err != nil {
			return out, false
		}
		nibbles := strings.Split(l, " ")
		out.key = make([]byte, len(nibbles))
		for i, n := range nibbles {
			tmp, _ := strconv.Atoi(n)
			out.key[i] = byte(tmp)
		}
	}

	{
		out.value, _ = file.ReadBytes('\n')
	}
	return out, true
}

func readInput(in io.Reader) []input {
	r := bufio.NewReader(in)
	var out []input
	i, ok := readLine(r)
	for ok {
		out = append(out, i)
		i, ok = readLine(r)
	}
	return out
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("%v", m.Alloc)
}

func main() {
	var t radix.KVStore
	if os.Args[1] == "trie" {
		t = trie.New()
	} else {
		t = radix.New()
	}

	dat := readInput(os.Stdin)

	st := time.Now()
	for _, d := range dat {
		t = t.Insert(d.key, d.value)
	}

	fmt.Println(time.Now().Sub(st).Seconds())
	dat = nil
	runtime.GC()
	PrintMemUsage()
	fmt.Println(t)
}
