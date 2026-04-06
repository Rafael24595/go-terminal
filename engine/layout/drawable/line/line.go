package line

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameLineDrawable = "LineDrawable"

type LineDrawable struct {
	loaded bool
	index  *IndexMeta
	lines  []text.Line
	cursor uint16
}

func NewLineDrawable(lines ...text.Line) *LineDrawable {
	return &LineDrawable{
		loaded: false,
		index:  &IndexMeta{},
		lines:  lines,
		cursor: 0,
	}
}

func LineDrawableFromLines(lines ...text.Line) drawable.Drawable {
	return NewLineDrawable(lines...).ToDrawable()
}

func (d *LineDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameLineDrawable,
		Code: "",
		Tags: make(set.Set[string]),
		Init: d.init,
		Draw: d.draw,
		Wipe: d.wipe,
	}
}

func (d *LineDrawable) init() {
	d.loaded = true

	d.index = computeIndexMeta(d.lines)
	d.cursor = 0
}

func (d *LineDrawable) wipe() {
	d.cursor = 0
}

func (d *LineDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

	if d.cursor >= uint16(len(d.lines)) {
		return make([]text.Line, 0), false
	}

	lines := indexLines(int(size.Cols), d.lines[d.cursor], d.index)
	d.cursor += 1

	return lines, d.cursor < uint16(len(d.lines))
}
