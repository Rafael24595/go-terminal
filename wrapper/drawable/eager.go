package wrapper_drawable

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type LinesEagerDrawable struct {
	rows  uint16
	cols  uint16
	lines []core.Line
}

func NewLinesEagerDrawable() *LinesEagerDrawable {
	return &LinesEagerDrawable{
		rows:  0,
		cols:  0,
		lines: make([]core.Line, 0),
	}
}

func (d *LinesEagerDrawable) init(size terminal.Winsize) {
	d.rows = size.Rows
	d.cols = size.Cols
}

func (d *LinesEagerDrawable) draw() ([]core.Line, bool) {
	lines := make([]core.Line, 0)
	for _, header := range d.lines {
		if int(d.cols) > header.Len() {
			lines = append(lines, header)
			continue
		}
		lines = append(lines, wrapLineWords(int(d.cols), header)...)
	}
	return lines, false
}

func (d *LinesEagerDrawable) ToDrawable() *drawable.Drawable {
	return &drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
