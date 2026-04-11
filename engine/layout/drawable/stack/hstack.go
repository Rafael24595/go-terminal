package stack

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/render/sink"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameHStackDrawable = "HStackDrawable"

const max_chunk = 100

type block struct {
	size  terminal.Winsize
	lines []text.Line
}

type HStackDrawable struct {
	loaded bool
	items  []layer
	fixed  []layer
}

func NewHStackDrawable(items ...drawable.Drawable) *HStackDrawable {
	layers := drawablesToLayer(items...)
	return &HStackDrawable{
		loaded: false,
		items:  layers,
		fixed:  make([]layer, 0),
	}
}

func HStackDrawableFromDrawables(items ...drawable.Drawable) drawable.Drawable {
	return NewVStackDrawable(items...).ToDrawable()
}

func (d *HStackDrawable) Unshift(items ...drawable.Drawable) *HStackDrawable {
	assert.False(d.loaded, err_new_elements)

	layers := drawablesToLayer(items...)
	d.items = append(layers, d.items...)

	return d
}

func (d *HStackDrawable) Push(items ...drawable.Drawable) *HStackDrawable {
	assert.False(d.loaded, err_new_elements)

	for _, item := range items {
		d.items = append(d.items,
			drawableToLayer(item),
		)
	}

	return d
}

func (d *HStackDrawable) UnshiftChunk(item drawable.Drawable, chunk uint16) *HStackDrawable {
	assert.False(d.loaded, err_new_elements)
	assert.True(chunk <= max_chunk, err_chunk_size, max_chunk)

	chunk = min(max_chunk, chunk)

	newLayer := layer{
		drawable: item,
		chunk:    chunk,
		status:   true,
	}

	d.items = append([]layer{newLayer}, d.items...)

	assert.LazyTrue(d.valideChunks, err_new_elements, max_chunk)

	return d
}

func (d *HStackDrawable) PushChunk(item drawable.Drawable, chunk uint16) *HStackDrawable {
	assert.False(d.loaded, err_new_elements)
	assert.True(chunk <= max_chunk, err_chunk_size, max_chunk)

	chunk = min(max_chunk, chunk)

	newLayer := layer{
		drawable: item,
		chunk:    chunk,
		status:   true,
	}

	d.items = append(d.items, newLayer)

	assert.LazyTrue(d.valideChunks, err_new_elements, max_chunk)

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

	d.fixed = d.items
	d.fixed = d.fixLayout()

	for i := range d.fixed {
		d.fixed[i].drawable.Init()
		d.fixed[i].status = true
	}
}

func (d *HStackDrawable) wipe() {
	d.fixed = d.items
	for i := range d.fixed {
		d.fixed[i].drawable.Wipe()
	}
}

func (d *HStackDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

	blocks, recalc := d.makeBlocks(size)
	lines := d.makeLines(blocks)

	if recalc {
		d.fixed = d.fixLayout()
	}

	return lines, len(d.fixed) > 0
}

func (d *HStackDrawable) makeBlocks(size terminal.Winsize) ([]block, bool) {
	buffer := make([]block, 0)
	recalcule := false

	for i := range d.fixed {
		if !d.fixed[i].status {
			continue
		}

		fixedSize := terminal.Winsize{
			Rows: size.Rows,
			Cols: (size.Cols * d.fixed[i].chunk) / 100,
		}

		lines, status := d.fixed[i].drawable.Draw(fixedSize)
		if !status {
			d.fixed[i].status = false
			recalcule = true
		}

		fixedLines := make([]text.Line, 0)
		for _, v := range lines {
			fixedLines = append(fixedLines,
				line.WrapLineWords(int(fixedSize.Cols), &v)...,
			)
		}

		buffer = append(buffer, block{
			size:  fixedSize,
			lines: fixedLines,
		})
	}

	return buffer, recalcule
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

func (d *HStackDrawable) fixLayout() []layer {
	chunks, zeroes := d.countChunks()

	assert.True(chunks <= max_chunk, err_new_elements, max_chunk)

	chunks = min(max_chunk, chunks)

	available := uint16(0)
	rest := uint16(0)
	if zeroes > 0 {
		remaining := math.SubClampZero(max_chunk, chunks)
		available = remaining / zeroes
		rest = remaining % zeroes
	}

	layers := make([]layer, 0, len(d.fixed))

	for _, v := range d.fixed {
		if !v.status {
			continue
		}

		chunk := min(max_chunk, v.chunk)
		if chunk == 0 {
			chunk = available
			if rest > 0 {
				chunk += 1
				rest -= 1
			}
		}

		layers = append(layers, layer{
			drawable: v.drawable,
			chunk:    chunk,
			status:   true,
		})
	}

	return layers
}

func (d *HStackDrawable) HasNext() bool {
	for _, item := range d.items {
		if item.status {
			return true
		}
	}
	return false
}

func (d *HStackDrawable) countChunks() (uint16, uint16) {
	chunks := uint16(0)
	zeroes := uint16(0)

	for _, i := range d.fixed {
		if !i.status {
			continue
		}

		chunks += i.chunk
		if i.chunk == 0 {
			zeroes += 1
		}
	}

	return chunks, zeroes
}

func (d *HStackDrawable) valideChunks() bool {
	chunks, _ := d.countChunks()
	return chunks <= max_chunk
}

func maxLines(buffer []block) int {
	colSize := 0
	for _, b := range buffer {
		colSize = max(colSize, len(b.lines))
	}
	return colSize
}
