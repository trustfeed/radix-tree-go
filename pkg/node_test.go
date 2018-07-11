package pkg

import (
	"testing"
)

func encodeString(s string) []byte {
	bts := []byte(s)
	out := make([]byte, len(bts)*2)
	for i, b := range bts {
		out[i*2] = b / 16
		out[i*2+1] = b % 16
	}
	return out
}

func TestInsert(t *testing.T) {
	vals := []struct{ k, v string }{
		{"dog", "one"},
		{"cat", "two"},
		{"doge", "three"},
		{"do", "four"},
	}

	var n node
	for _, v := range vals {
		n = insert(n, encodeString(v.k), []byte(v.v))
	}

	for _, v := range vals {
		str := string(lookup(n, encodeString(v.k)))
		if str != v.v {
			t.Errorf("looked up data differs; expected %s got %s", v.v, str)
		}
	}
}
