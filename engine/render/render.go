package render

import (
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

type Adapter func([]text.Line, winsize.Winsize) string
type RawAdapter func([]text.Line, winsize.Winsize) []string

type Render struct {
	Render Adapter
}

type RenderBuilder struct {
	render Adapter
}

func NewBuilder(render Adapter) *RenderBuilder {
	return &RenderBuilder{
		render: render,
	}
}

func (b *RenderBuilder) ToRender() Render {
	return Render{
		Render: b.render,
	}
}
