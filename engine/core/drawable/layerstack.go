package drawable

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type LayerStack struct {
	items  []Drawable
	status map[int]bool
}

func (d *LayerStack) Init(size terminal.Winsize) {
	for _, item := range d.items {
		item.Init(size )
	}
}

func (d *LayerStack) Push(item Drawable) {
	d.items = append(d.items, item)
	d.status[len(d.items)] = true
}

func (d *LayerStack) Draw() ([]core.Line, bool) {
	buffer := make([]core.Line, 0)
	gStatus := false

	for i, item := range d.items {
		if !d.status[i] {
			continue
		}

		lines, status := item.Draw()
		if !status {
			d.status[i] = false
		}

		buffer = append(buffer, lines...)
		gStatus = status || gStatus

		if gStatus {
			break
		}
	}

	return buffer, gStatus
}
