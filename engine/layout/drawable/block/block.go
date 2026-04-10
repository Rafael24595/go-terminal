package block

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameBlockDrawable = "BlockDrawable"

type BlockDrawable struct {
	loaded   bool
	drawable drawable.Drawable
}

func NewBlockDrawable(drawable drawable.Drawable) *BlockDrawable {
	return &BlockDrawable{
		loaded:   false,
		drawable: drawable,
	}
}

func BlockDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewBlockDrawable(drawable).ToDrawable()
}

func BlockDrawableFromLines(lines ...text.Line) drawable.Drawable {
	return BlockDrawableFromDrawable(
		line.NewLineDrawable(lines...).ToDrawable(),
	)
}

func BlockDrawableFromString(txt ...string) drawable.Drawable {
	lines := text.LineFromFragments(
		text.FragmentsFromString(txt...)...,
	)
	return BlockDrawableFromLines(*lines)
}

func (d *BlockDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameBlockDrawable,
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

func (d *BlockDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, "the drawable should be initialized before draw")

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
