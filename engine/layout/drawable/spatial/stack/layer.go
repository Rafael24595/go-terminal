package stack

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
)

type layer[T math.Number] struct {
	drawable drawable.Drawable
	chunk    chunk.Chunk[T]
	value    T
	status   bool
}

func layerFromLayer[T math.Number](other layer[T], value T) layer[T] {
	return layer[T]{
		drawable: other.drawable,
		chunk:    other.chunk,
		value:    value,
		status:   true,
	}
}

func layerFromDrawable[T math.Number](item drawable.Drawable, value chunk.Chunk[T], cols T) layer[T] {
	return layer[T]{
		drawable: item,
		chunk:    value,
		value:    cols,
		status:   true,
	}
}

func layersFromDrawables[T math.Number](chk chunk.Chunk[T], value T, items ...drawable.Drawable) []layer[T] {
	layers := make([]layer[T], len(items))
	for i, item := range items {
		layers[i] = layerFromDrawable(
			item, chk, value,
		)
	}
	return layers
}
