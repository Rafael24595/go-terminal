package draw

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type DrawContext struct {
	State *state.UIState
	Size  terminal.Winsize
}

func NewDrawContext(stt *state.UIState, size terminal.Winsize) *DrawContext {
	return &DrawContext{
		State: stt,
		Size:  size,
	}
}
