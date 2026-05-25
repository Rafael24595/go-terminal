package stack

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
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

func (u *VStackUnit) Unshift(units ...drawable.Unit) *VStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(
		chunk.Dynamic[winsize.Rows](), 0, units...,
	)

	u.items = append(layers, u.items...)
	return u
}

func (u *VStackUnit) Push(units ...drawable.Unit) *VStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	for _, unit := range units {
		u.items = append(u.items,
			layerFromUnit(chunk.Dynamic[winsize.Rows](), 0, unit),
		)
	}

	return u
}

func (u *VStackUnit) UnshiftChunk(
	unit drawable.Unit,
	chunk chunk.Chunk[winsize.Rows],
) *VStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(chunk, 0, unit)

	u.items = append(layers, u.items...)
	return u
}

func (u *VStackUnit) PushChunk(unit drawable.Unit, chunk chunk.Chunk[winsize.Rows]) *VStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	u.items = append(u.items,
		layerFromUnit(chunk, 0, unit),
	)

	return u
}

func (u *VStackUnit) Size() uint {
	return uint(len(u.items))
}

func (u *VStackUnit) Units() []drawable.Unit {
	units := make([]drawable.Unit, len(u.items))
	for i := range u.items {
		units[i] = u.items[i].unit
	}
	return units
}

func (u *VStackUnit) ToUnit() drawable.Unit {
	if u.isAnemic() {
		unit := u.items[0].unit
		return unit.AddTag(AnemicStack)
	}

	return drawable.NewBuilder().
		Name(NameVStack).
		MergeTags(u.tags()).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *VStackUnit) isAnemic() bool {
	if len(u.items) != 1 {
		return false
	}
	return u.items[0].chunk.IsAnemic()
}

func (u *VStackUnit) tags() set.Set[string] {
	tags := set.NewSet[string]()
	for i := range u.items {
		tags.Merge(u.items[i].unit.Tags)
	}
	return tags
}

func (u *VStackUnit) init() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *VStackUnit) lazyInit(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	u.fixed = u.items
	u.fixed = u.fixLayout(size)

	for i := range u.fixed {
		u.fixed[i].unit.Drawable.Init()
		u.fixed[i].status = true
	}
}

func (u *VStackUnit) wipe() {
	u.lazyLoaded = false

	u.fixed = u.items
	for i := range u.fixed {
		u.fixed[i].unit.Drawable.Wipe()
	}
}

func (u *VStackUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	u.lazyInit(size)

	if !u.size.Eq(size) {
		u.fixed = u.fixLayout(size)
		u.size = size
	}

	lines, recalc := u.makeLines(size)

	if !u.size.Eq(size) || recalc {
		u.fixed = u.fixLayout(size)
	}

	return lines, len(u.fixed) > 0
}

func (u *VStackUnit) makeLines(size winsize.Winsize) ([]text.Line, bool) {
	buffer := make([]text.Line, 0, size.Rows)
	recalcule := false

	for i := range u.fixed {
		if !u.fixed[i].status {
			continue
		}

		bufferLen := winsize.Rows(len(buffer))
		remaining := size.Rows.Sub(bufferLen)
		if remaining == 0 {
			break
		}

		rows := remaining
		if u.fixed[i].chunk.Sized {
			value := u.fixed[i].value
			rows = min(value, remaining)
		}

		fixedSize := winsize.New(rows, size.Cols)

		lines, status := drain.Unit(fixedSize, u.fixed[i].unit, true)
		if !status {
			u.fixed[i].status = false
			recalcule = true
		}

		linesLen := winsize.Rows(len(lines))
		if linesLen < rows && u.fixed[i].chunk.Sized {
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

func (u *VStackUnit) fixLayout(size winsize.Winsize) []layer[winsize.Rows] {
	layers := make([]layer[winsize.Rows], 0, len(u.fixed))

	for _, v := range u.fixed {
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

func (u *VStackUnit) HasNext() bool {
	items := u.items
	if u.lazyLoaded {
		items = u.fixed
	}

	for _, item := range items {
		if item.status {
			return true
		}
	}

	return false
}
