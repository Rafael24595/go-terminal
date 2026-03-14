package set

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](size ...int) Set[T] {
	s := 0

	if len(size) > 0 {
		s = size[0]
	}

	return make(Set[T], s)
}

func SetFrom[T comparable](slice ...T) Set[T] {
	s := NewSet[T](len(slice))

	for _, v := range slice {
		s.Add(v)
	}

	return s
}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Any(other Set[T]) bool {
	a, b := s, other
	if len(a) > len(b) {
		a, b = b, a
	}

	for k := range a {
		if _, ok := b[k]; ok {
			return true
		}
	}

	return false
}
