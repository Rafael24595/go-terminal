package padding

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

type SizeHint[T math.Number] struct {
	size T
	fill bool
}

func Fixed[T math.Number](cols T) SizeHint[T] {
	return SizeHint[T]{
		size: cols,
	}
}

func Maximize[T math.Number]() SizeHint[T] {
	return SizeHint[T]{
		fill: true,
	}
}

func (h SizeHint[T]) Min(maxSize T) T {
	if h.fill {
		return maxSize
	}
	return min(h.size, maxSize)
}
