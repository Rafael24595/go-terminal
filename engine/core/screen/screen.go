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
	IgnoreParents bool
	Screen        *Screen
	Pager         state.PagerContext
}

type Definition struct {
	RequireKeys []key.Key
}

func DefinitionFromKeys(keys ...key.Key) Definition {
	return Definition{
		RequireKeys: keys,
	}
}

func ScreenResultFromScreen(screen *Screen) ScreenResult {
	return ScreenResult{
		IgnoreParents: false,
		Screen:        screen,
		Pager:         state.PagerContext{},
	}
}

func ScreenResultFromUIState(stt *state.UIState) ScreenResult {
	return ScreenResult{
		IgnoreParents: false,
		Screen:        nil,
		Pager:         stt.Pager,
	}
}

func EmptyScreenResult() ScreenResult {
	return ScreenResult{
		IgnoreParents: false,
		Screen:        nil,
		Pager:         state.PagerContext{},
	}
}

type Screen struct {
	//Init func (ctx)
	Name       func() string
	Definition func() Definition
	Update     func(*state.UIState, ScreenEvent) ScreenResult
	View       func(state.UIState) core.ViewModel
}
