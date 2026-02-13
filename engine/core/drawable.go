package core

import "github.com/Rafael24595/go-terminal/engine/terminal"

type Drawable[T any] struct {
	Init func(size terminal.Winsize) *Drawable[T]
	Draw func() ([]Line, bool)
}
