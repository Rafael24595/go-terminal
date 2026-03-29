package block

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameBlockDrawable = "BlockDrawable"

type BlockDrawable struct {
	initialized bool
	size        terminal.Winsize
	drawable    drawable.Drawable
}

func NewBlockDrawable(drawable drawable.Drawable) *BlockDrawable {
	return &BlockDrawable{
		initialized: false,
		size:        terminal.Winsize{},
		drawable:    drawable,
	}
}

func BlockDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewBlockDrawable(drawable).ToDrawable()
}

func (d *BlockDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameBlockDrawable,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *BlockDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.size = size

	d.drawable.Init(size)
}

func (d *BlockDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	lines := make([]text.Line, 0)
	for range d.size.Rows {
		lns, hasNext := d.drawable.Draw()
		lines = append(lines, lns...)

		if !hasNext {
			return lines, hasNext
		}
	}

	return lines, true
}
