package stack

import (
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/model/chunk"
)

type layer struct {
	drawable drawable.Drawable
	status   bool
}

func layersFromDrawables(items ...drawable.Drawable) []layer {
	layers := make([]layer, len(items))
	for i, item := range items {
		layers[i] = layerFromDrawable(item)
	}
	return layers
}

func layerFromDrawable(item drawable.Drawable) layer {
	return layer{
		drawable: item,
		status:   true,
	}
}

type chunkLayer struct {
	drawable drawable.Drawable
	chunk    chunk.Chunk
	cols     uint16
	status   bool
}

func chunkLayerFromLayer(other chunkLayer, cols uint16) chunkLayer {
	return chunkLayer{
		drawable: other.drawable,
		chunk:    other.chunk,
		cols:     cols,
		status:   true,
	}
}

func chunkLayerFromDrawable(item drawable.Drawable, chunk chunk.Chunk, cols uint16) chunkLayer {
	return chunkLayer{
		drawable: item,
		chunk:    chunk,
		cols:     cols,
		status:   true,
	}
}

func chunkLayersFromDrawables(chk chunk.Chunk, cols uint16, items ...drawable.Drawable) []chunkLayer {
	layers := make([]chunkLayer, len(items))
	for i, item := range items {
		layers[i] = chunkLayerFromDrawable(
			item, chk, cols,
		)
	}
	return layers
}
