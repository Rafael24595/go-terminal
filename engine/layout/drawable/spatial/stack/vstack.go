package stack

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameVStack = "vstack_unit"

type VStackUnit struct {
	loaded     bool
	lazyLoaded bool
	size       winsize.Winsize
	items      []layer[winsize.Rows]
	fixed      []layer[winsize.Rows]
}

func NewVStack(units ...drawable.Unit) *VStackUnit {
	layers := layersFromUnits(
		chunk.Dynamic[winsize.Rows](), 0, units...,
	)

	return &VStackUnit{
		loaded:     false,
		lazyLoaded: false,
		size:       winsize.Winsize{},
		items:      layers,
		fixed:      make([]layer[winsize.Rows], 0),
	}
}

func VStackFromUnits(units ...drawable.Unit) drawable.Unit {
	return NewVStack(units...).ToUnit()
}

func (d *VStackUnit) Unshift(units ...drawable.Unit) *VStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(
		chunk.Dynamic[winsize.Rows](), 0, units...,
	)

	d.items = append(layers, d.items...)
	return d
}

func (d *VStackUnit) Push(units ...drawable.Unit) *VStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	for _, unit := range units {
		d.items = append(d.items,
			layerFromUnit(chunk.Dynamic[winsize.Rows](), 0, unit),
		)
	}

	return d
}

func (d *VStackUnit) UnshiftChunk(
	unit drawable.Unit,
	chunk chunk.Chunk[winsize.Rows],
) *VStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(chunk, 0, unit)

	d.items = append(layers, d.items...)
	return d
}

func (d *VStackUnit) PushChunk(unit drawable.Unit, chunk chunk.Chunk[winsize.Rows]) *VStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	d.items = append(d.items,
		layerFromUnit(chunk, 0, unit),
	)

	return d
}

func (d *VStackUnit) Size() uint {
	return uint(len(d.items))
}

func (d *VStackUnit) Units() []drawable.Unit {
	units := make([]drawable.Unit, len(d.items))
	for i := range d.items {
		units[i] = d.items[i].unit
	}
	return units
}

func (d *VStackUnit) ToUnit() drawable.Unit {
	if d.isAnemic() {
		unit := d.items[0].unit
		return unit.AddTag(AnemicStack)
	}

	return drawable.NewBuilder().
		Name(NameVStack).
		MergeTags(d.tags()).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *VStackUnit) isAnemic() bool {
	if len(d.items) != 1 {
		return false
	}
	return d.items[0].chunk.IsAnemic()
}

func (d *VStackUnit) tags() set.Set[string] {
	tags := set.NewSet[string]()
	for i := range d.items {
		tags.Merge(d.items[i].unit.Tags)
	}
	return tags
}

func (d *VStackUnit) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *VStackUnit) lazyInit(size winsize.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	d.fixed = d.items
	d.fixed = d.fixLayout(size)

	for i := range d.fixed {
		d.fixed[i].unit.Drawable.Init()
		d.fixed[i].status = true
	}
}

func (d *VStackUnit) wipe() {
	d.lazyLoaded = false

	d.fixed = d.items
	for i := range d.fixed {
		d.fixed[i].unit.Drawable.Wipe()
	}
}

func (d *VStackUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
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

func (d *VStackUnit) makeLines(size winsize.Winsize) ([]text.Line, bool) {
	buffer := make([]text.Line, 0, size.Rows)
	recalcule := false

	for i := range d.fixed {
		if !d.fixed[i].status {
			continue
		}

		bufferLen := winsize.Rows(len(buffer))
		remaining := size.Rows.Sub(bufferLen)
		if remaining == 0 {
			break
		}

		rows := remaining
		if d.fixed[i].chunk.Sized {
			value := d.fixed[i].value
			rows = min(value, remaining)
		}

		fixedSize := winsize.New(rows, size.Cols)

		lines, status := drain.Unit(fixedSize, d.fixed[i].unit, true)
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

func (d *VStackUnit) fixLayout(size winsize.Winsize) []layer[winsize.Rows] {
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

func (d *VStackUnit) HasNext() bool {
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
