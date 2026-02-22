package line

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type LazyDrawable struct {
	initialized bool
	rows        uint16
	cols        uint16
	meta        *IndexMeta
	lines       []core.Line
	cursor      uint16
}

func NewLazyDrawable(lines ...core.Line) *LazyDrawable {
	return &LazyDrawable{
		initialized: false,
		rows:        0,
		cols:        0,
		meta:        &IndexMeta{},
		lines:       lines,
		cursor:      0,
	}
}

func LazyDrawableFromLines(lines ...core.Line) core.Drawable {
	return NewLazyDrawable(lines...).ToDrawable()
}

func (d *LazyDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.rows = size.Rows
	d.cols = size.Cols
	d.meta = computeIndexMeta(d.lines)
	d.cursor = 0
}

func (d *LazyDrawable) draw() ([]core.Line, bool) {
	assert.AssertTrue(d.initialized, "the drawable should be initialized before draw")

	if d.cursor >= uint16(len(d.lines)) {
		return make([]core.Line, 0), false
	}

	lines := indexLines(int(d.cols), d.lines[d.cursor], d.meta)
	d.cursor += 1

	return lines, d.cursor < uint16(len(d.lines))
}

func (d *LazyDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
