package line

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "line_drawable"

type LineDrawable struct {
	loaded bool
	index  *indexMeta
	lines  []text.Line
	source []text.Line
}

func New(lines ...text.Line) *LineDrawable {
	return &LineDrawable{
		loaded: false,
		index:  &indexMeta{},
		lines:  lines,
	}
}

func DrawableFromLines(lines ...text.Line) drawable.Drawable {
	return New(lines...).ToDrawable()
}

func (d *LineDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
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

func (d *LineDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

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
