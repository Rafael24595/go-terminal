package stack

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
)

type layer[T math.Number] struct {
	unit   drawable.Unit
	chunk  chunk.Chunk[T]
	value  T
	status bool
}

func layerFromLayer[T math.Number](other layer[T], value T) layer[T] {
	return layer[T]{
		unit:   other.unit,
		chunk:  other.chunk,
		value:  value,
		status: true,
	}
}

func layerFromUnit[T math.Number](
	chunk chunk.Chunk[T],
	value T,
	unit drawable.Unit,
) layer[T] {
	return layer[T]{
		unit:   unit,
		chunk:  chunk,
		value:  value,
		status: true,
	}
}

func layersFromUnits[T math.Number](
	chunk chunk.Chunk[T],
	value T,
	units ...drawable.Unit,
) []layer[T] {
	layers := make([]layer[T], len(units))
	for i := range units {
		layers[i] = layerFromUnit(
			chunk, value, units[i],
		)
	}
	return layers
}
