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
	Pager         state.PagerState
	Cursor        state.CursorState
}

type Definition struct {
	RequireKeys []key.Key
}

func NewScreenResult(screen *Screen, pager state.PagerState, cursor state.CursorState) ScreenResult {
	return ScreenResult{
		IgnoreParents: false,
		Screen:        screen,
		Cursor:        cursor,
		Pager:         pager,
	}
}

func ScreenResultFromScreen(screen *Screen) ScreenResult {
	return ScreenResult{
		IgnoreParents: false,
		Screen:        screen,
		Pager:         state.EmptyPagerState(),
		Cursor:        state.EmptyCursorState(),
	}
}

func ScreenResultFromUIState(state state.UIState) ScreenResult {
	return ScreenResult{
		IgnoreParents: false,
		Screen:        nil,
		Pager:         state.Pager,
		Cursor:        state.Cursor,
	}
}

func EmptyScreenResult() ScreenResult {
	return ScreenResult{
		IgnoreParents: false,
		Screen:        nil,
		Pager:         state.EmptyPagerState(),
		Cursor:        state.EmptyCursorState(),
	}
}

type Screen struct {
	//Init func (ctx)
	Name       func() string
	Definition func() Definition
	Update     func(state.UIState, ScreenEvent) ScreenResult
	View       func(state.UIState) core.ViewModel
}
