package layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type apply func(state *state.UIState, vm viewmodel.ViewModel, size terminal.Winsize) []text.Line

type Layout struct {
	Apply apply
}

func NewLayout(apply apply) Layout {
	return Layout{
		Apply: apply,
	}
}
