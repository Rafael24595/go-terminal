package line

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type LinesEagerDrawable struct {
	rows   uint16
	cols   uint16
	status bool
	lines  []core.Line
}

func NewLinesEagerDrawable(lines ...core.Line) *LinesEagerDrawable {
	return &LinesEagerDrawable{
		rows:   0,
		cols:   0,
		status: true,
		lines:  lines,
	}
}

func LinesEagerDrawableFromLines(lines ...core.Line) core.Drawable {
	return NewLinesEagerDrawable(lines...).ToDrawable()
}

func (d *LinesEagerDrawable) init(size terminal.Winsize) {
	d.rows = size.Rows
	d.cols = size.Cols
	d.status = true
}

func (d *LinesEagerDrawable) draw() ([]core.Line, bool) {
	lines := make([]core.Line, 0)

	if !d.status {
		return lines, d.status
	}

	for _, header := range d.lines {
		if int(d.cols) > header.Len() {
			lines = append(lines, header)
			continue
		}
		lines = append(lines, WrapLineWords(int(d.cols), header)...)
	}

	d.status = false
	return lines, d.status
}

func (d *LinesEagerDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
