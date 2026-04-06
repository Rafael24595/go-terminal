package stack

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameStackDrawable = "StackDrawable"

type layer struct {
	drawable drawable.Drawable
	status   bool
}

type StackDrawable struct {
	loaded bool
	items  []layer
}

func NewStackDrawable(items ...drawable.Drawable) *StackDrawable {
	layers := drawableToLayer(items...)
	return &StackDrawable{
		loaded: false,
		items:  layers,
	}
}

func StackDrawableFromDrawables(items ...drawable.Drawable) drawable.Drawable {
	return NewStackDrawable(items...).ToDrawable()
}

func (d *StackDrawable) Unshift(items ...drawable.Drawable) *StackDrawable {
	assert.False(d.loaded, "no new elements should be added after initialization")

	layers := drawableToLayer(items...)
	d.items = append(layers, d.items...)
	return d
}

func (d *StackDrawable) Push(items ...drawable.Drawable) *StackDrawable {
	assert.False(d.loaded, "no new elements should be added after initialization")

	for _, item := range items {
		d.items = append(d.items, layer{
			drawable: item,
			status:   true,
		})
	}
	return d
}

func (d *StackDrawable) Size() uint {
	return uint(len(d.items))
}

func (d *StackDrawable) Take(code string) (drawable.Drawable, bool) {
	for i, v := range d.items {
		if v.drawable.Code == code {
			target := v.drawable
			d.items = append(d.items[:i], d.items[i+1:]...)
			return target, true
		}
	}
	return drawable.Drawable{}, false
}

func (d *StackDrawable) Items() []drawable.Drawable {
	items := make([]drawable.Drawable, len(d.items))
	for i := range d.items {
		items[i] = d.items[i].drawable
	}
	return items
}

func (d *StackDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameStackDrawable,
		Code: d.Code(),
		Tags: d.Tags(),
		Init: func() {
			d.Init()
		},
		Wipe: func() {
			d.Wipe()
		},
		Draw: d.Draw,
	}
}

func (d *StackDrawable) Code() string {
	var sb strings.Builder
	for i := range d.items {
		_, _ = sb.Write([]byte(d.items[i].drawable.Code))
	}
	return sb.String()
}

func (d *StackDrawable) Tags() set.Set[string] {
	tags := set.NewSet[string]()
	for i := range d.items {
		tags.Merge(d.items[i].drawable.Tags)
	}
	return tags
}

func (d *StackDrawable) Init() *StackDrawable {
	d.loaded = true

	for i := range d.items {
		d.items[i].drawable.Init()
		d.items[i].status = true
	}
	return d
}

func (d *StackDrawable) Wipe() *StackDrawable {
	for i := range d.items {
		d.items[i].drawable.Wipe()
	}
	return d
}

func (d *StackDrawable) Draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

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

func (d *StackDrawable) HasNext() bool {
	for _, item := range d.items {
		if item.status {
			return true
		}
	}
	return false
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
