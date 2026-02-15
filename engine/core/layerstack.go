package core

import (
	"iter"

	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type layer struct {
	drawable Drawable
	status   bool
}

type LayerStack struct {
	items []layer
}

func NewLayerStack() *LayerStack {
	return &LayerStack{
		items: make([]layer, 0),
	}
}

func (d *LayerStack) Init(size terminal.Winsize) *LayerStack {
	for i := range d.items {
		d.items[i].drawable.Init(size)
		d.items[i].status = true
	}
	return d
}

func (d *LayerStack) Unshift(items ...Drawable) *LayerStack {
	newLayers := make([]layer, len(items))
	for i, item := range items {
		newLayers[i] = layer{
			drawable: item,
			status:   true,
		}
	}
	d.items = append(newLayers, d.items...)
	return d
}

func (d *LayerStack) Shift(items ...Drawable) *LayerStack {
	for _, item := range items {
		d.items = append(d.items, layer{
			drawable: item,
			status:   true,
		})
	}
	return d
}

func (d *LayerStack) Draw() ([]Line, bool) {
	buffer := make([]Line, 0)
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

func (d *LayerStack) Iterator() iter.Seq[[]Line] {
	return func(yield func([]Line) bool) {
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

func (d *LayerStack) HasNext() bool {
	for _, item := range d.items {
		if item.status {
			return true
		}
	}
	return false
}
