package stack

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/sink"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const NameHStack = "hstack_unit"

type block struct {
	size  winsize.Winsize
	lines []text.Line
}

type HStackUnit struct {
	loaded     bool
	lazyLoaded bool
	size       winsize.Winsize
	items      []layer[winsize.Cols]
	fixed      []layer[winsize.Cols]
}

func NewHStack(units ...drawable.Unit) *HStackUnit {
	layers := layersFromUnits(
		chunk.Dynamic[winsize.Cols](), 0, units...,
	)

	return &HStackUnit{
		loaded:     false,
		lazyLoaded: false,
		size:       winsize.Winsize{},
		items:      layers,
		fixed:      make([]layer[winsize.Cols], 0),
	}
}

func HStackFromUnits(units ...drawable.Unit) drawable.Unit {
	return NewHStack(units...).ToUnit()
}

func (d *HStackUnit) Unshift(units ...drawable.Unit) *HStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(
		chunk.Dynamic[winsize.Cols](), 0, units...,
	)

	d.items = append(layers, d.items...)
	return d
}

func (d *HStackUnit) Push(units ...drawable.Unit) *HStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	for _, unit := range units {
		d.items = append(d.items,
			layerFromUnit(chunk.Dynamic[winsize.Cols](), 0, unit),
		)
	}

	return d
}

func (d *HStackUnit) UnshiftChunk(unit drawable.Unit, chunk chunk.Chunk[winsize.Cols]) *HStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(chunk, 0, unit)

	d.items = append(layers, d.items...)
	return d
}

func (d *HStackUnit) PushChunk(unit drawable.Unit, chunk chunk.Chunk[winsize.Cols]) *HStackUnit {
	assert.False(d.loaded, drawable.MessageNewElement)

	d.items = append(d.items,
		layerFromUnit(chunk, 0, unit),
	)
	return d
}

func (d *HStackUnit) Size() uint {
	return uint(len(d.items))
}

func (d *HStackUnit) Units() []drawable.Unit {
	units := make([]drawable.Unit, len(d.items))
	for i := range d.items {
		units[i] = d.items[i].unit
	}
	return units
}

func (d *HStackUnit) ToUnit() drawable.Unit {
	if d.isAnemic() {
		unit := d.items[0].unit
		return unit.AddTag(AnemicStack)
	}

	return drawable.NewBuilder().
		Name(NameHStack).
		MergeTags(d.tags()).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *HStackUnit) isAnemic() bool {
	if len(d.items) != 1 {
		return false
	}
	return d.items[0].chunk.IsAnemic()
}

func (d *HStackUnit) tags() set.Set[string] {
	tags := set.NewSet[string]()
	for i := range d.items {
		tags.Merge(d.items[i].unit.Tags)
	}
	return tags
}

func (d *HStackUnit) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *HStackUnit) lazyInit(size winsize.Winsize) {
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

func (d *HStackUnit) wipe() {
	d.lazyLoaded = false

	d.fixed = d.items
	for i := range d.fixed {
		d.fixed[i].unit.Drawable.Wipe()
	}
}

func (d *HStackUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	if !d.size.Eq(size) {
		d.fixed = d.fixLayout(size)
		d.size = size
	}

	blocks, recalc := d.makeBlocks(size)
	lines := d.makeLines(blocks)

	if !d.size.Eq(size) || recalc {
		d.fixed = d.fixLayout(size)
	}

	return lines, len(d.fixed) > 0
}

func (d *HStackUnit) makeBlocks(size winsize.Winsize) ([]block, bool) {
	buffer := make([]block, len(d.fixed))
	recalcule := false

	maxHeight := 0

	canGrow := make([]bool, len(d.fixed))
	for i := range d.fixed {
		canGrow[i] = d.fixed[i].status
	}

	for {
		didGrow := false

		for i := range d.fixed {
			if !d.fixed[i].status || (maxHeight > 0 && len(buffer[i].lines) >= maxHeight) {
				continue
			}

			inheritCols := d.inheritCols(size, buffer, i)
			fixedSize := winsize.Winsize{
				Rows: size.Rows,
				Cols: d.fixed[i].value + inheritCols,
			}

			drawable := d.fixed[i].unit.Drawable
			lines, status := drawable.Draw(fixedSize)
			if !status {
				d.fixed[i].status = false
				canGrow[i] = false
				recalcule = true
			}

			if len(lines) == 0 {
				continue
			}

			wrapped := make([]text.Line, 0)
			for _, v := range lines {
				wrapped = append(wrapped,
					wrap.Line(fixedSize.Cols, &v)...,
				)
			}

			buffer[i].size = fixedSize
			buffer[i].lines = append(buffer[i].lines, wrapped...)

			if len(buffer[i].lines) > maxHeight {
				maxHeight = len(buffer[i].lines)
			}

			didGrow = true
		}

		if !didGrow {
			break
		}

		shouldContinue := false
		for i := range d.fixed {
			if d.fixed[i].status && len(buffer[i].lines) < maxHeight {
				shouldContinue = true
				break
			}
		}

		if !shouldContinue {
			break
		}
	}

	return buffer, recalcule
}

func (d *HStackUnit) inheritCols(
	size winsize.Winsize,
	buffer []block,
	bufferIndex int,
) winsize.Cols {
	if bufferIndex == 0 {
		return 0
	}

	block := buffer[bufferIndex-1]
	lineIndex := len(buffer[bufferIndex].lines)

	if len(block.lines) <= lineIndex {
		return 0
	}

	line := block.lines[lineIndex]
	if text.FragmentMeasure(size.Cols, line.Text...) != 0 {
		return 0
	}

	return d.fixed[bufferIndex-1].value
}

func (d *HStackUnit) makeLines(blocks []block) []text.Line {
	buffer := make([]text.Line, 0)
	for i := range maxLines(blocks) {
		line := text.EmptyLine()
		for _, b := range blocks {
			if i >= len(b.lines) {
				continue
			}

			l := b.lines[i]
			result := sink.ApplySinks(&l, b.size.Cols)

			line.CopyMeta(result)
			line.PushFragments(result.Text...)
		}
		buffer = append(buffer, *line)
	}

	return buffer
}

func (d *HStackUnit) fixLayout(size winsize.Winsize) []layer[winsize.Cols] {
	layers := make([]layer[winsize.Cols], 0, len(d.fixed))
	available, rest := d.calcSpace(size)

	for _, v := range d.fixed {
		if !v.status {
			continue
		}

		chk := v.chunk

		chunk := winsize.Cols(0)
		if chk.Sized {
			chunk = min(size.Cols, chk.Adapter(size.Cols))
		} else {
			chunk = available
			if rest > 0 {
				chunk += 1
				rest -= 1
			}
		}

		layers = append(layers,
			layerFromLayer(v, chunk),
		)
	}

	assert.LazyTrue(func() bool {
		return d.valideChunks(size)
	}, drawable.MessageNewElement, size.Cols)

	return layers
}

func (d *HStackUnit) calcSpace(size winsize.Winsize) (winsize.Cols, winsize.Cols) {
	cols, zeroes := d.countCols(size)

	assert.True(cols <= size.Cols, drawable.MessageNewElement, size.Cols)

	if zeroes == 0 {
		return 0, 0
	}

	cols = min(size.Cols, cols)
	remaining := size.Cols.Sub(cols)

	cZeroes := winsize.Cols(zeroes)

	available := remaining / cZeroes
	rest := remaining % cZeroes

	return available, rest
}

func (d *HStackUnit) HasNext() bool {
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

func (d *HStackUnit) countCols(size winsize.Winsize) (winsize.Cols, uint16) {
	cols := winsize.Cols(0)
	zeroes := uint16(0)

	for _, i := range d.fixed {
		if !i.status {
			continue
		}

		chk := i.chunk
		if !chk.Sized {
			zeroes += 1
		} else {
			cols += chk.Adapter(size.Cols)
		}
	}

	return cols, zeroes
}

func (d *HStackUnit) valideChunks(size winsize.Winsize) bool {
	cols, _ := d.countCols(size)
	return cols <= size.Cols
}

func maxLines(buffer []block) int {
	colSize := 0
	for _, b := range buffer {
		colSize = max(colSize, len(b.lines))
	}
	return colSize
}
