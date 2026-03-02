package line

import (
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type EagerLoopDrawable struct {
	initialized bool
	eager       *EagerDrawable
}

func NewEagerLoopDrawable(lines ...text.Line) *EagerLoopDrawable {
	return &EagerLoopDrawable{
		initialized: false,
		eager:       NewEagerDrawable(lines...),
	}
}

func EagerLoopDrawableFromLines(lines ...text.Line) drawable.Drawable {
	return NewEagerLoopDrawable(lines...).ToDrawable()
}

func (d *EagerLoopDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.eager.init(size)
}

func (d *EagerLoopDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines, status := d.eager.draw()
	if !status {
		d.eager.status = true
	}

	return lines, true
}

func (d *EagerLoopDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
