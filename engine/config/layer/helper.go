package layer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

func Fixed[T math.Number](chk T) Option[T] {
	return WithChunk(chunk.Fixed(chk))
}

func Percent[T math.Number](chk T) Option[T] {
	return WithChunk(chunk.Percent(chk))
}
