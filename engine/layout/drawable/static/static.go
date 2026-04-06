package static

import (
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameStaticDrawable = "StaticDrawable"

type StaticDrawable struct {
	drawable drawable.Drawable
}

func NewStaticDrawable(drawable drawable.Drawable) *StaticDrawable {
	return &StaticDrawable{
		drawable: drawable,
	}
}

func StaticDrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return NewStaticDrawable(drawable).ToDrawable()
}

func (s *StaticDrawable) init() {}

func (s *StaticDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	return s.drawable.Draw(size)
}

func (d *StaticDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameStaticDrawable,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.drawable.Init,
		Draw: d.draw,
	}
}
