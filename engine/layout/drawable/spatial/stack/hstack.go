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

func (u *HStackUnit) Unshift(units ...drawable.Unit) *HStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(
		chunk.Dynamic[winsize.Cols](), 0, units...,
	)

	u.items = append(layers, u.items...)
	return u
}

func (u *HStackUnit) Push(units ...drawable.Unit) *HStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	for _, unit := range units {
		u.items = append(u.items,
			layerFromUnit(chunk.Dynamic[winsize.Cols](), 0, unit),
		)
	}

	return u
}

func (u *HStackUnit) UnshiftChunk(unit drawable.Unit, chunk chunk.Chunk[winsize.Cols]) *HStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	layers := layersFromUnits(chunk, 0, unit)

	u.items = append(layers, u.items...)
	return u
}

func (u *HStackUnit) PushChunk(unit drawable.Unit, chunk chunk.Chunk[winsize.Cols]) *HStackUnit {
	assert.False(u.loaded, drawable.MessageNewElement)

	u.items = append(u.items,
		layerFromUnit(chunk, 0, unit),
	)
	return u
}

func (u *HStackUnit) Size() uint {
	return uint(len(u.items))
}

func (u *HStackUnit) Units() []drawable.Unit {
	units := make([]drawable.Unit, len(u.items))
	for i := range u.items {
		units[i] = u.items[i].unit
	}
	return units
}

func (u *HStackUnit) ToUnit() drawable.Unit {
	if u.isAnemic() {
		unit := u.items[0].unit
		return unit.AddTag(AnemicStack)
	}

	return drawable.NewBuilder().
		Name(NameHStack).
		MergeTags(u.tags()).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *HStackUnit) isAnemic() bool {
	if len(u.items) != 1 {
		return false
	}
	return u.items[0].chunk.IsAnemic()
}

func (u *HStackUnit) tags() set.Set[string] {
	tags := set.NewSet[string]()
	for i := range u.items {
		tags.Merge(u.items[i].unit.Tags)
	}
	return tags
}

func (u *HStackUnit) init() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *HStackUnit) lazyInit(size winsize.Winsize) {
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

func (u *HStackUnit) wipe() {
	u.lazyLoaded = false

	u.fixed = u.items
	for i := range u.fixed {
		u.fixed[i].unit.Drawable.Wipe()
	}
}

func (u *HStackUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	u.lazyInit(size)

	if !u.size.Eq(size) {
		u.fixed = u.fixLayout(size)
		u.size = size
	}

	blocks, recalc := u.makeBlocks(size)
	lines := u.makeLines(blocks)

	if !u.size.Eq(size) || recalc {
		u.fixed = u.fixLayout(size)
	}

	return lines, len(u.fixed) > 0
}

func (u *HStackUnit) makeBlocks(size winsize.Winsize) ([]block, bool) {
	buffer := make([]block, len(u.fixed))
	recalcule := false

	maxHeight := 0

	canGrow := make([]bool, len(u.fixed))
	for i := range u.fixed {
		canGrow[i] = u.fixed[i].status
	}

	for {
		didGrow := false

		for i := range u.fixed {
			if !u.fixed[i].status || (maxHeight > 0 && len(buffer[i].lines) >= maxHeight) {
				continue
			}

			inheritCols := u.inheritCols(size, buffer, i)
			fixedSize := winsize.Winsize{
				Rows: size.Rows,
				Cols: u.fixed[i].value + inheritCols,
			}

			drawable := u.fixed[i].unit.Drawable
			lines, status := drawable.Draw(fixedSize)
			if !status {
				u.fixed[i].status = false
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
		for i := range u.fixed {
			if u.fixed[i].status && len(buffer[i].lines) < maxHeight {
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

func (u *HStackUnit) inheritCols(
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

	return u.fixed[bufferIndex-1].value
}

func (u *HStackUnit) makeLines(blocks []block) []text.Line {
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

func (u *HStackUnit) fixLayout(size winsize.Winsize) []layer[winsize.Cols] {
	layers := make([]layer[winsize.Cols], 0, len(u.fixed))
	available, rest := u.calcSpace(size)

	for _, v := range u.fixed {
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
		return u.valideChunks(size)
	}, drawable.MessageNewElement, size.Cols)

	return layers
}

func (u *HStackUnit) calcSpace(size winsize.Winsize) (winsize.Cols, winsize.Cols) {
	cols, zeroes := u.countCols(size)

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

func (u *HStackUnit) HasNext() bool {
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

func (u *HStackUnit) countCols(size winsize.Winsize) (winsize.Cols, uint16) {
	cols := winsize.Cols(0)
	zeroes := uint16(0)

	for _, i := range u.fixed {
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

func (u *HStackUnit) valideChunks(size winsize.Winsize) bool {
	cols, _ := u.countCols(size)
	return cols <= size.Cols
}

func maxLines(buffer []block) int {
	colSize := 0
	for _, b := range buffer {
		colSize = max(colSize, len(b.lines))
	}
	return colSize
}
