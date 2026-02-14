package drawable

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type Drawable struct {
	Init func(size terminal.Winsize)
	Draw func() ([]core.Line, bool)
}
