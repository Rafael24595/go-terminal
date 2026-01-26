package core

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type apply func(state *state.UIState, vm ViewModel, size terminal.Winsize) []Line

type Layout struct {
	Apply apply
}

func NewLayout(apply apply) Layout {
	return Layout{
		Apply: apply,
	}
}
