package line

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type EagerDrawable struct {
	initialized bool
	rows        uint16
	cols        uint16
	status      bool
	lines       []core.Line
}

func NewEagerDrawable(lines ...core.Line) *EagerDrawable {
	return &EagerDrawable{
		initialized: false,
		rows:        0,
		cols:        0,
		status:      true,
		lines:       lines,
	}
}

func EagerDrawableFromLines(lines ...core.Line) core.Drawable {
	return NewEagerDrawable(lines...).ToDrawable()
}

func (d *EagerDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.rows = size.Rows
	d.cols = size.Cols
	d.status = true
}

func (d *EagerDrawable) draw() ([]core.Line, bool) {
	assert.AssertTrue(d.initialized, "the drawable should be initialized before draw")

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

func (d *EagerDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
