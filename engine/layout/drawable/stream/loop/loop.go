package loop

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

const NameLoopDrawable = "LoopDrawable"

type LoopDrawable struct {
	loaded   bool
	drawable drawable.Drawable
}

func NewLoopDrawable(drawable drawable.Drawable) *LoopDrawable {
	return &LoopDrawable{
		loaded:   false,
		drawable: drawable,
	}
}

func LoopDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewLoopDrawable(drawable).ToDrawable()
}

func (d *LoopDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameLoopDrawable,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.drawable.Wipe,
		Draw: d.draw,
	}
}

func (d *LoopDrawable) init() {
	d.loaded = true

	d.drawable.Init()
}

func (d *LoopDrawable) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	lines, status := d.drawable.Draw(size)
	if !status {
		d.drawable.Wipe()
	}

	return lines, true
}
