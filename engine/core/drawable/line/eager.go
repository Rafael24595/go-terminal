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

	for _, line := range d.lines {
		if int(d.cols) >= line.Len() {
			lines = append(lines, line)
			continue
		}
		lines = append(lines, WrapLineWords(int(d.cols), line)...)
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
