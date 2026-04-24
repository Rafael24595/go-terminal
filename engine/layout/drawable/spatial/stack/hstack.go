package stack

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/sink"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameHStackDrawable = "HStackDrawable"

type block struct {
	size  winsize.Winsize
	lines []text.Line
}

type HStackDrawable struct {
	loaded     bool
	lazyLoaded bool
	size       winsize.Winsize
	items      []chunkLayer
	fixed      []chunkLayer
}

func NewHStackDrawable(items ...drawable.Drawable) *HStackDrawable {
	layers := chunkLayersFromDrawables(chunk.Dynamic(), 0, items...)
	return &HStackDrawable{
		loaded:     false,
		lazyLoaded: false,
		size:       winsize.Winsize{},
		items:      layers,
		fixed:      make([]chunkLayer, 0),
	}
}

func HStackDrawableFromDrawables(items ...drawable.Drawable) drawable.Drawable {
	return NewVStackDrawable(items...).ToDrawable()
}

func (d *HStackDrawable) Unshift(items ...drawable.Drawable) *HStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	layers := chunkLayersFromDrawables(chunk.Dynamic(), 0, items...)
	d.items = append(layers, d.items...)

	return d
}

func (d *HStackDrawable) Push(items ...drawable.Drawable) *HStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	for _, item := range items {
		d.items = append(d.items,
			chunkLayerFromDrawable(item, chunk.Dynamic(), 0),
		)
	}

	return d
}

func (d *HStackDrawable) UnshiftChunk(item drawable.Drawable, chunk chunk.Chunk) *HStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	newLayer := chunkLayerFromDrawable(item, chunk, 0)

	d.items = append([]chunkLayer{newLayer}, d.items...)

	return d
}

func (d *HStackDrawable) PushChunk(item drawable.Drawable, chunk chunk.Chunk) *HStackDrawable {
	assert.False(d.loaded, drawable.MessageNewElement)

	newLayer := chunkLayerFromDrawable(item, chunk, 0)

	d.items = append(d.items, newLayer)

	return d
}

func (d *HStackDrawable) Size() uint {
	return uint(len(d.items))
}

func (d *HStackDrawable) Take(code string) (drawable.Drawable, bool) {
	for i, v := range d.items {
		if v.drawable.Code == code {
			target := v.drawable
			d.items = append(d.items[:i], d.items[i+1:]...)
			return target, true
		}
	}
	return drawable.Drawable{}, false
}

func (d *HStackDrawable) Items() []drawable.Drawable {
	items := make([]drawable.Drawable, len(d.items))
	for i := range d.items {
		items[i] = d.items[i].drawable
	}
	return items
}

func (d *HStackDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameVStackDrawable,
		Code: d.code(),
		Tags: d.tags(),
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *HStackDrawable) code() string {
	var sb strings.Builder
	for i := range d.items {
		_, _ = sb.Write([]byte(d.items[i].drawable.Code))
	}
	return sb.String()
}

func (d *HStackDrawable) tags() set.Set[string] {
	tags := set.NewSet[string]()
	for i := range d.items {
		tags.Merge(d.items[i].drawable.Tags)
	}
	return tags
}

func (d *HStackDrawable) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *HStackDrawable) lazyInit(size winsize.Winsize) {
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

func (d *HStackDrawable) wipe() {
	d.lazyLoaded = false

	d.fixed = d.items
	for i := range d.fixed {
		d.fixed[i].drawable.Wipe()
	}
}

func (d *HStackDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
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

func (d *HStackDrawable) makeBlocks(size winsize.Winsize) ([]block, bool) {
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
				Cols: d.fixed[i].cols + inheritCols,
			}

			lines, status := d.fixed[i].drawable.Draw(fixedSize)
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
				wrapped = append(wrapped, line.WrapLineWords(int(fixedSize.Cols), &v)...)
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

func (d *HStackDrawable) inheritCols(
	size winsize.Winsize,
	buffer []block,
	bufferIndex int,
) uint16 {
	if bufferIndex == 0 {
		return 0
	}

	block := buffer[bufferIndex-1]
	lineIndex := len(buffer[bufferIndex].lines)

	if len(block.lines) <= lineIndex {
		return 0
	}

	line := block.lines[lineIndex]
	if text.FragmentMeasure(int(size.Cols), line.Text...) != 0 {
		return 0
	}

	return d.fixed[bufferIndex-1].cols
}

func (d *HStackDrawable) makeLines(blocks []block) []text.Line {
	buffer := make([]text.Line, 0)
	for i := range maxLines(blocks) {
		line := text.EmptyLine()
		for _, b := range blocks {
			if i >= len(b.lines) {
				continue
			}

			l := b.lines[i]
			result := sink.ApplySinks(&l, int(b.size.Cols))

			line.CopyMeta(result)
			line.PushFragments(result.Text...)
		}
		buffer = append(buffer, *line)
	}

	return buffer
}

func (d *HStackDrawable) fixLayout(size winsize.Winsize) []chunkLayer {
	layers := make([]chunkLayer, 0, len(d.fixed))
	available, rest := d.calcSpace(size)

	for _, v := range d.fixed {
		if !v.status {
			continue
		}

		chk := v.chunk

		chunk := uint16(0)
		if chk.Sized {
			chunk = min(size.Cols, chk.Adapter(size))
		} else {
			chunk = available
			if rest > 0 {
				chunk += 1
				rest -= 1
			}
		}

		layers = append(layers,
			chunkLayerFromLayer(v, chunk),
		)
	}

	assert.LazyTrue(func() bool {
		return d.valideChunks(size)
	}, drawable.MessageNewElement, size.Cols)

	return layers
}

func (d *HStackDrawable) calcSpace(size winsize.Winsize) (uint16, uint16) {
	cols, zeroes := d.countCols(size)

	assert.True(cols <= size.Cols, drawable.MessageNewElement, size.Cols)

	if zeroes == 0 {
		return 0, 0
	}

	cols = min(size.Cols, cols)
	remaining := math.SubClampZero(size.Cols, cols)

	available := remaining / zeroes
	rest := remaining % zeroes

	return available, rest
}

func (d *HStackDrawable) HasNext() bool {
	for _, item := range d.items {
		if item.status {
			return true
		}
	}
	return false
}

func (d *HStackDrawable) countCols(size winsize.Winsize) (uint16, uint16) {
	cols := uint16(0)
	zeroes := uint16(0)

	for _, i := range d.fixed {
		if !i.status {
			continue
		}

		chk := i.chunk
		if !chk.Sized {
			zeroes += 1
		} else {
			cols += chk.Adapter(size)
		}
	}

	return cols, zeroes
}

func (d *HStackDrawable) valideChunks(size winsize.Winsize) bool {
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
