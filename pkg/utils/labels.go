package utils

// CheckSubset returns true if map "a" is subset of map "b"
func CheckSubset(a, b map[string]string) bool {
	for k, v := range a {
		if bv, ok := b[k]; !ok {
			return false
		} else if v != bv {
			return false
		}
	}

	return true
}
