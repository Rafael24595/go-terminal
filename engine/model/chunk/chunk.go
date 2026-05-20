package chunk

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

const max_chunk = 100

const (
	err_chunk_size = "chunk value should be less or equals than %s"
)

type chunkAdapter[T math.Number] func(size T) T

type Chunk[T math.Number] struct {
	isAnemic bool
	Adapter  chunkAdapter[T]
	Sized    bool
}

func Dynamic[T math.Number]() Chunk[T] {
	return Chunk[T]{
		isAnemic: true,
		Adapter:  fixAdapter(T(0)),
		Sized:    false,
	}
}

func Fixed[T math.Number](fix T) Chunk[T] {
	return Chunk[T]{
		isAnemic: false,
		Adapter:  fixAdapter(fix),
		Sized:    true,
	}
}

func Percent[T math.Number](chk T) Chunk[T] {
	if chk > max_chunk {
		assert.Unreachable(err_chunk_size, max_chunk)
		chk = max_chunk
	}

	return Chunk[T]{
		isAnemic: false,
		Adapter:  perAdapter(chk),
		Sized:    true,
	}
}

func (c Chunk[T]) IsAnemic() bool {
	return c.isAnemic
}

func perAdapter[T math.Number](chunk T) chunkAdapter[T] {
	return func(size T) T {
		return (size * chunk) / 100
	}
}

func fixAdapter[T math.Number](cols T) chunkAdapter[T] {
	return func(size T) T {
		if cols > size {
			return size
		}
		return cols
	}
}
