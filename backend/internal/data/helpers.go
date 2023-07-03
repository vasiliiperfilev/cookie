package data

func ReverseMap[K comparable, V comparable](m map[K]V) map[V]K {
	n := make(map[V]K, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

func EqualArrays[T comparable](a1 []T, a2 []T) bool {
	m := map[T]bool{}
	for _, el := range a1 {
		m[el] = true
	}
	for _, el := range a2 {
		if _, ok := m[el]; ok {
			delete(m, el)
		} else {
			return false
		}
	}
	if len(m) > 0 {
		return false
	} else {
		return true
	}
}
