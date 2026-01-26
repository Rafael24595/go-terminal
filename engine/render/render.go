package render

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type render func([]core.Line, terminal.Winsize) string

type Render struct {
	Render render
}

func NewRender(render render) Render {
	return Render{
		Render: render,
	}
}