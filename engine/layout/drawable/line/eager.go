package line

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type EagerDrawable struct {
	initialized bool
	size        terminal.Winsize
	status      bool
	lines       []text.Line
}

func NewEagerDrawable(lines ...text.Line) *EagerDrawable {
	return &EagerDrawable{
		initialized: false,
		size:        terminal.Winsize{},
		status:      true,
		lines:       lines,
	}
}

func EagerDrawableFromLines(lines ...text.Line) drawable.Drawable {
	return NewEagerDrawable(lines...).ToDrawable()
}

func EagerDrawableFromString(txt ...string) drawable.Drawable {
	lines := text.LineFromFragments(
		text.FragmentsFromString(txt...)...,
	)
	return EagerDrawableFromLines(lines)
}

func (d *EagerDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.size = size
	d.status = true
}

func (d *EagerDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines := make([]text.Line, 0)

	if !d.status {
		return lines, d.status
	}

	for _, line := range d.lines {
		lines = append(lines, WrapLineWords(int(d.size.Cols), line)...)
	}

	d.status = false
	return lines, d.status
}

func (d *EagerDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
