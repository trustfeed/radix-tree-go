package trie

import (
	"testing"
)

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
