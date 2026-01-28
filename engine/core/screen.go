package core

import "github.com/Rafael24595/go-terminal/engine/app/state"

type ScreenEvent struct {
	Key string
}

type ScreenResult struct {
	State         state.UIState
	IgnoreParents bool
	Screen        *Screen
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
	Name   func() string
	Update func(state.UIState, ScreenEvent) ScreenResult
	View   func(state.UIState) ViewModel
}
