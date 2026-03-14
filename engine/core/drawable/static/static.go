package static

import (
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/text"
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

func (s *StaticDrawable) init(size terminal.Winsize) {}

func (s *StaticDrawable) draw() ([]text.Line, bool) {
	return s.drawable.Draw()
}

func (d *StaticDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameStaticDrawable,
		Init: d.init,
		Draw: d.draw,
	}
}
