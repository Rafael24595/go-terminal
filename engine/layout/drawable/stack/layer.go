package stack

import "github.com/Rafael24595/go-terminal/engine/layout/drawable"

type layer struct {
	drawable drawable.Drawable
	status   bool
}

func drawableToLayer(items ...drawable.Drawable) []layer {
	layers := make([]layer, len(items))
	for i, item := range items {
		layers[i] = layer{
			drawable: item,
			status:   true,
		}
	}
	return layers
}
