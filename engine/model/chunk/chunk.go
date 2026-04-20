package chunk

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
)

const max_chunk = 100

const (
	err_chunk_size = "chunk value should be less or equals than %s"
)

type chunkAdapter func(size winsize.Winsize) uint16

type Chunk struct {
	Adapter chunkAdapter
	Sized   bool
}

func Dynamic() Chunk {
	return Chunk{
		Adapter: colsAdapter(0),
		Sized:   false,
	}
}

func Colums(cols uint16) Chunk {
	return Chunk{
		Adapter: colsAdapter(cols),
		Sized:   true,
	}
}

func Percent(chk uint16) Chunk {
	if chk > max_chunk {
		assert.Unreachable(err_chunk_size, max_chunk)
		chk = max_chunk
	}

	return Chunk{
		Adapter: percAdapter(chk),
		Sized:   true,
	}
}

func percAdapter(chunk uint16) chunkAdapter {
	return func(size winsize.Winsize) uint16 {
		return (size.Cols * chunk) / 100
	}
}

func colsAdapter(cols uint16) chunkAdapter {
	return func(size winsize.Winsize) uint16 {
		if cols > size.Cols {
			return size.Cols
		}
		return cols
	}
}
