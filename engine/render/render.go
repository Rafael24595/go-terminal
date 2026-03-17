package render

import (
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type render func([]text.Line, terminal.Winsize) string

type Render struct {
	Render render
}

func NewRender(render render) Render {
	return Render{
		Render: render,
	}
}
