package stack

import (
	"iter"

	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type layer struct {
	drawable drawable.Drawable
	status   bool
}

type StackDrawable struct {
	initialized bool
	items       []layer
}

func NewStackDrawable(items ...drawable.Drawable) *StackDrawable {
	layers := drawableToLayer(items...)
	return &StackDrawable{
		initialized: false,
		items:       layers,
	}
}

func StackDrawableFromDrawables(items ...drawable.Drawable) drawable.Drawable {
	return NewStackDrawable(items...).ToDrawable()
}

func (d *StackDrawable) Init(size terminal.Winsize) *StackDrawable {
	d.initialized = true
	
	for i := range d.items {
		d.items[i].drawable.Init(size)
		d.items[i].status = true
	}
	return d
}

func (d *StackDrawable) Unshift(items ...drawable.Drawable) *StackDrawable {
	assert.False(d.initialized, "no new elements should be added after initialization")

	layers := drawableToLayer(items...)
	d.items = append(layers, d.items...)
	return d
}

func (d *StackDrawable) Shift(items ...drawable.Drawable) *StackDrawable {
	assert.False(d.initialized, "no new elements should be added after initialization")

	for _, item := range items {
		d.items = append(d.items, layer{
			drawable: item,
			status:   true,
		})
	}
	return d
}

func (d *StackDrawable) Draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	buffer := make([]text.Line, 0)
	gStatus := false

	for i := range d.items {
		if !d.items[i].status {
			continue
		}

		lines, status := d.items[i].drawable.Draw()
		if !status {
			d.items[i].status = false
		}

		buffer = append(buffer, lines...)
		gStatus = status || gStatus

		if gStatus {
			break
		}
	}

	return buffer, gStatus
}

func (d *StackDrawable) Iterator() iter.Seq[[]text.Line] {
	return func(yield func([]text.Line) bool) {
		for {
			lines, content := d.Draw()
			if !yield(lines) {
				return
			}

			if !content {
				return
			}
		}
	}
}

func (d *StackDrawable) HasNext() bool {
	for _, item := range d.items {
		if item.status {
			return true
		}
	}
	return false
}

func (d *StackDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: func(size terminal.Winsize) {
			d.Init(size)
		},
		Draw: d.Draw,
	}
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
