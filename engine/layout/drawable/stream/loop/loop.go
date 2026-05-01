package loop

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "loop_drawable"

type LoopDrawable struct {
	loaded   bool
	drawable drawable.Drawable
}

func New(drawable drawable.Drawable) *LoopDrawable {
	return &LoopDrawable{
		loaded:   false,
		drawable: drawable,
	}
}

func DrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return New(drawable).ToDrawable()
}

func (d *LoopDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: Name,
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
