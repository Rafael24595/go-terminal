package render

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type Render struct {
	Render func(core.ViewModel, terminal.Winsize) string
}
