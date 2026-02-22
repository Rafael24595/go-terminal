package line

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type EagerLoopDrawable struct {
	initialized bool
	eager       *EagerDrawable
}

func NewEagerLoopDrawable(lines ...core.Line) *EagerLoopDrawable {
	return &EagerLoopDrawable{
		initialized: false,
		eager:       NewEagerDrawable(lines...),
	}
}

func EagerLoopDrawableFromLines(lines ...core.Line) core.Drawable {
	return NewEagerLoopDrawable(lines...).ToDrawable()
}

func (d *EagerLoopDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.eager.init(size)
}

func (d *EagerLoopDrawable) draw() ([]core.Line, bool) {
	assert.AssertTrue(d.initialized, "the drawable should be initialized before draw")

	lines, status := d.eager.draw()
	if !status {
		d.eager.status = true
	}

	return lines, true
}

func (d *EagerLoopDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
