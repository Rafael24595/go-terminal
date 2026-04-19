package stack

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameVStackDrawable = "VStackDrawable"

type VStackDrawable struct {
	loaded bool
	items  []layer
}

func NewVStackDrawable(items ...drawable.Drawable) *VStackDrawable {
	layers := layersFromDrawables(items...)
	return &VStackDrawable{
		loaded: false,
		items:  layers,
	}
}

func VStackDrawableFromDrawables(items ...drawable.Drawable) drawable.Drawable {
	return NewVStackDrawable(items...).ToDrawable()
}

func (d *VStackDrawable) Unshift(items ...drawable.Drawable) *VStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	layers := layersFromDrawables(items...)
	d.items = append(layers, d.items...)
	return d
}

func (d *VStackDrawable) Push(items ...drawable.Drawable) *VStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	for _, item := range items {
		d.items = append(d.items, layer{
			drawable: item,
			status:   true,
		})
	}
	return d
}

func (d *VStackDrawable) Size() uint {
	return uint(len(d.items))
}

func (d *VStackDrawable) Take(code string) (drawable.Drawable, bool) {
	for i, v := range d.items {
		if v.drawable.Code == code {
			target := v.drawable
			d.items = append(d.items[:i], d.items[i+1:]...)
			return target, true
		}
	}
	return drawable.Drawable{}, false
}

func (d *VStackDrawable) Items() []drawable.Drawable {
	items := make([]drawable.Drawable, len(d.items))
	for i := range d.items {
		items[i] = d.items[i].drawable
	}
	return items
}

func (d *VStackDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameVStackDrawable,
		Code: d.code(),
		Tags: d.tags(),
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *VStackDrawable) code() string {
	var sb strings.Builder
	for i := range d.items {
		_, _ = sb.Write([]byte(d.items[i].drawable.Code))
	}
	return sb.String()
}

func (d *VStackDrawable) tags() set.Set[string] {
	tags := set.NewSet[string]()
	for i := range d.items {
		tags.Merge(d.items[i].drawable.Tags)
	}
	return tags
}

func (d *VStackDrawable) init() {
	d.loaded = true

	for i := range d.items {
		d.items[i].drawable.Init()
		d.items[i].status = true
	}
}

func (d *VStackDrawable) wipe() {
	for i := range d.items {
		d.items[i].drawable.Wipe()
	}
}

func (d *VStackDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	buffer := make([]text.Line, 0)
	gStatus := false

	for i := range d.items {
		if !d.items[i].status {
			continue
		}

		lines, status := d.items[i].drawable.Draw(size)
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

func (d *VStackDrawable) HasNext() bool {
	for _, item := range d.items {
		if item.status {
			return true
		}
	}
	return false
}
