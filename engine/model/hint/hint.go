package hint

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

type Size[T math.Number] struct {
	size T
	fill bool
}

func Fixed[T math.Number](cols T) Size[T] {
	return Size[T]{
		size: cols,
	}
}

func Maximize[T math.Number]() Size[T] {
	return Size[T]{
		fill: true,
	}
}

func (h Size[T]) Min(maxSize T) T {
	if h.fill {
		return maxSize
	}
	return min(h.size, maxSize)
}
