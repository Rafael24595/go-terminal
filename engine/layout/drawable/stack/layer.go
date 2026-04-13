package stack

import "github.com/Rafael24595/go-terminal/engine/layout/drawable"

type layer struct {
	drawable drawable.Drawable
	chunk    uint16
	sized    bool
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
	return chunkLayerFromDrawable(item, 0)
}

func chunkLayerFromLayer(other layer, chunk uint16) layer {
	return layer{
		drawable: other.drawable,
		chunk:    chunk,
		sized:    other.sized,
		status:   true,
	}
}

func chunkLayerFromDrawable(item drawable.Drawable, chunk uint16) layer {
	return layer{
		drawable: item,
		chunk:    chunk,
		sized:    chunk != 0,
		status:   true,
	}
}
