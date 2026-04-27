package stack

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameVStackDrawable = "VStackDrawable"

type VStackDrawable struct {
	loaded     bool
	lazyLoaded bool
	size       winsize.Winsize
	items      []layer[winsize.Rows]
	fixed      []layer[winsize.Rows]
}

func NewVStackDrawable(items ...drawable.Drawable) *VStackDrawable {
	layers := layersFromDrawables(chunk.Dynamic[winsize.Rows](), 0, items...)
	return &VStackDrawable{
		loaded:     false,
		lazyLoaded: false,
		size:       winsize.Winsize{},
		items:      layers,
		fixed:      make([]layer[winsize.Rows], 0),
	}
}

func VStackDrawableFromDrawables(items ...drawable.Drawable) drawable.Drawable {
	return NewVStackDrawable(items...).ToDrawable()
}

func (d *VStackDrawable) Unshift(items ...drawable.Drawable) *VStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	layers := layersFromDrawables(chunk.Dynamic[winsize.Rows](), 0, items...)
	d.items = append(layers, d.items...)
	return d
}

func (d *VStackDrawable) Push(items ...drawable.Drawable) *VStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	for _, item := range items {
		d.items = append(d.items,
			layerFromDrawable(item, chunk.Dynamic[winsize.Rows](), 0),
		)
	}

	return d
}

func (d *VStackDrawable) UnshiftChunk(item drawable.Drawable, chunk chunk.Chunk[winsize.Rows]) *VStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	newLayer := layerFromDrawable(item, chunk, 0)

	d.items = append([]layer[winsize.Rows]{newLayer}, d.items...)

	return d
}

func (d *VStackDrawable) PushChunk(item drawable.Drawable, chunk chunk.Chunk[winsize.Rows]) *VStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	newLayer := layerFromDrawable(item, chunk, 0)

	d.items = append(d.items, newLayer)

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
	d.lazyLoaded = false
}

func (d *VStackDrawable) lazyInit(size winsize.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	d.fixed = d.items
	d.fixed = d.fixLayout(size)

	for i := range d.fixed {
		d.fixed[i].drawable.Init()
		d.fixed[i].status = true
	}
}

func (d *VStackDrawable) wipe() {
	for i := range d.items {
		d.items[i].drawable.Wipe()
	}
}

func (d *VStackDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	if !d.size.Eq(size) {
		d.fixed = d.fixLayout(size)
		d.size = size
	}

	lines, recalc := d.makeLines(size)

	if !d.size.Eq(size) || recalc {
		d.fixed = d.fixLayout(size)
	}

	return lines, len(d.fixed) > 0
}

func (d *VStackDrawable) makeLines(size winsize.Winsize) ([]text.Line, bool) {
	buffer := make([]text.Line, 0, size.Rows)
	recalcule := false

	for i := range d.fixed {
		if !d.fixed[i].status {
			continue
		}

		bufferLen := winsize.Rows(len(buffer))
		remaining := math.SubClampZero(size.Rows, bufferLen)
		if remaining <= 0 {
			break
		}

		rows := remaining
		if d.fixed[i].chunk.Sized {
			value := d.fixed[i].value
			rows = min(value, remaining)
		}

		fixedSize := winsize.Winsize{
			Rows: winsize.Rows(rows),
			Cols: size.Cols,
		}

		lines, status := d.fixed[i].drawable.Draw(fixedSize)
		if !status {
			d.fixed[i].status = false
			recalcule = true
		}

		linesLen := winsize.Rows(len(lines))
		if linesLen < rows && d.fixed[i].chunk.Sized {
			padded := make([]text.Line, rows)
			copy(padded, lines)
			lines = padded
			linesLen = rows
		}

		limit := min(rows, linesLen)
		buffer = append(buffer, lines[:limit]...)
	}

	return buffer, recalcule
}

func (d *VStackDrawable) fixLayout(size winsize.Winsize) []layer[winsize.Rows] {
	layers := make([]layer[winsize.Rows], 0, len(d.fixed))

	for _, v := range d.fixed {
		if !v.status {
			continue
		}

		chk := v.chunk

		chunk := winsize.Rows(0)
		if chk.Sized {
			chunk = min(size.Rows, chk.Adapter(size.Rows))
		}

		layers = append(layers,
			layerFromLayer(v, chunk),
		)
	}

	return layers
}

func (d *VStackDrawable) HasNext() bool {
	items := d.items
	if d.lazyLoaded {
		items = d.fixed
	}

	for _, item := range items {
		if item.status {
			return true
		}
	}

	return false
}
