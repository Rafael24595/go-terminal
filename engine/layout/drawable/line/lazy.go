package line

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type LazyDrawable struct {
	initialized bool
	rows        uint16
	cols        uint16
	index       *IndexMeta
	lines       []text.Line
	cursor      uint16
}

func NewLazyDrawable(lines ...text.Line) *LazyDrawable {
	return &LazyDrawable{
		initialized: false,
		rows:        0,
		cols:        0,
		index:       &IndexMeta{},
		lines:       lines,
		cursor:      0,
	}
}

func LazyDrawableFromLines(lines ...text.Line) drawable.Drawable {
	return NewLazyDrawable(lines...).ToDrawable()
}

func (d *LazyDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.rows = size.Rows
	d.cols = size.Cols
	d.index = computeIndexMeta(d.lines)
	d.cursor = 0
}

func (d *LazyDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	if d.cursor >= uint16(len(d.lines)) {
		return make([]text.Line, 0), false
	}

	lines := indexLines(int(d.cols), d.lines[d.cursor], d.index)
	d.cursor += 1

	return lines, d.cursor < uint16(len(d.lines))
}

func (d *LazyDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
