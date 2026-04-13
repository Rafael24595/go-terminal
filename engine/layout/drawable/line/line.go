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
	source []text.Line
}

func NewLineDrawable(lines ...text.Line) *LineDrawable {
	return &LineDrawable{
		loaded: false,
		index:  &IndexMeta{},
		lines:  lines,
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

	d.lines = TokenizeLines(d.lines...)
	d.source = text.CloneLines(d.lines...)

	d.index = computeIndexMeta(d.lines)
}

func (d *LineDrawable) wipe() {
	d.source = d.lines
}

func (d *LineDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

	if len(d.source) == 0 {
		return make([]text.Line, 0), false
	}

	cursor, remain := WrapNextLine(size.Cols, d.source, d.index)
	d.source = remain

	result := make([]text.Line, 0)
	if cursor != nil {
		result = append(result, *cursor)
	}

	return result, len(d.source) > 0
}
