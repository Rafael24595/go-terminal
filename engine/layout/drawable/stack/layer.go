package stack

import "github.com/Rafael24595/go-terminal/engine/layout/drawable"

type layer struct {
	drawable drawable.Drawable
	chunk    uint16
	status   bool
}

func drawablesToLayer(items ...drawable.Drawable) []layer {
	layers := make([]layer, len(items))
	for i, item := range items {
		layers[i] = drawableToLayer(item)
	}
	return layers
}

func drawableToLayer(item drawable.Drawable) layer {
	return layer{
		drawable: item,
		chunk:    0,
		status:   true,
	}
}
