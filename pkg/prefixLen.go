package pkg

// Find the length of a shared prefix between the two byte slices.
func prefixLen(x, y []byte) int {
	out := 0
	for out < len(x) && out < len(y) {
		if x[out] != y[out] {
			return out
		}
		out++
	}
	return out
}
