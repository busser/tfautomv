package slices

// Contains reports whether v is within s.
func Contains[T comparable](s []T, v T) bool {
	return Index(s, v) >= 0
}

// Index returns the index of the first instance of v in s, or -1 if v is not
// present in s.
func Index[T comparable](s []T, v T) int {
	for i := range s {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// Equal returns whether a and b's contents are identical.
func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
