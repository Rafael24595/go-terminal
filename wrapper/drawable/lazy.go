package wrapper_drawable

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type LinesLazyDrawable struct {
	rows   uint16
	cols   uint16
	meta   *IndexMeta
	lines  []core.Line
	cursor uint16
}

func NewLinesLazyDrawable() *LinesLazyDrawable {
	return &LinesLazyDrawable{
		rows:   0,
		cols:   0,
		meta:   &IndexMeta{},
		lines:  make([]core.Line, 0),
		cursor: 0,
	}
}

func (d *LinesLazyDrawable) init(size terminal.Winsize) {
	d.rows = size.Rows
	d.cols = size.Cols
	d.meta = computeIndexMeta(d.lines)
	d.cursor = 0
}

func (d *LinesLazyDrawable) draw() ([]core.Line, bool) {
	if d.cursor >= uint16(len(d.lines)) {
		return make([]core.Line, 0), false
	}

	lines := indexLines(int(d.cols), d.lines[d.cursor], d.meta)
	d.cursor += 1
	return lines, d.cursor < uint16(len(d.lines))
}

func (d *LinesLazyDrawable) ToDrawable() *drawable.Drawable {
	return &drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
