package drawable

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type layer struct {
	drawable Drawable
	status   bool
}

type LayerStack struct {
	items []layer
}

func (d *LayerStack) Init(size terminal.Winsize) {
	for _, item := range d.items {
		item.drawable.Init(size)
	}
}

func (d *LayerStack) Unshift(items ...Drawable) {
	newLayers := make([]layer, len(items))
	for i, item := range items {
		newLayers[i] = layer{
			drawable: item,
			status:   true,
		}
	}
	d.items = append(newLayers, d.items...)
}

func (d *LayerStack) Shift(items ...Drawable) {
	for _, item := range items {
		d.items = append(d.items, layer{
			drawable: item,
			status:   true,
		})
	}
}

func (d *LayerStack) Draw() ([]core.Line, bool) {
	buffer := make([]core.Line, 0)
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
