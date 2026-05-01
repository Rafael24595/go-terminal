package block

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "block_drawable"

type BlockDrawable struct {
	loaded   bool
	drawable drawable.Drawable
}

func New(drawable drawable.Drawable) *BlockDrawable {
	return &BlockDrawable{
		loaded:   false,
		drawable: drawable,
	}
}

func DrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return New(drawable).ToDrawable()
}

func DrawableFromLines(lines ...text.Line) drawable.Drawable {
	return DrawableFromDrawable(
		line.New(lines...).ToDrawable(),
	)
}

func DrawableFromString(txt ...string) drawable.Drawable {
	lines := text.LineFromFragments(
		text.FragmentsFromString(txt...)...,
	)
	return DrawableFromLines(*lines)
}

func (d *BlockDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.drawable.Wipe,
		Draw: d.draw,
	}
}

func (d *BlockDrawable) init() {
	d.loaded = true

	d.drawable.Init()
}

func (d *BlockDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	lines := make([]text.Line, 0)
	for range size.Rows {
		lns, hasNext := d.drawable.Draw(size)
		lines = append(lines, lns...)

		if !hasNext {
			return lines, hasNext
		}
	}

	return lines, true
}
