package screen

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
)

type ScreenEvent struct {
	Key key.Key
}

type ScreenResult struct {
	State         state.UIState
	IgnoreParents bool
	Screen        *Screen
}

type Definition struct {
	RequireKeys []key.Key
}

func NewScreenResult(state state.UIState, screen *Screen) ScreenResult {
	return ScreenResult{
		State:  state,
		Screen: screen,
	}
}

func ScreenResultFromState(state state.UIState) ScreenResult {
	return ScreenResult{
		State: state,
	}
}

type Screen struct {
	//Init func (ctx)
	Name       func() string
	Definition func() Definition
	Update     func(state.UIState, ScreenEvent) ScreenResult
	View       func(state.UIState) core.ViewModel
}
