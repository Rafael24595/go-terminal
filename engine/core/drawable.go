package core

import (
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type Drawable struct {
	Init func(size terminal.Winsize)
	Draw func() ([]Line, bool)
}
