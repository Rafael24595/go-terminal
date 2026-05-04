package screen

import "github.com/Rafael24595/go-reacterm-core/engine/app/state"

type Result struct {
	IgnoreParents bool
	Screen        *Screen
	Pager         state.PagerContext
}

func ResultFromScreen(screen *Screen) Result {
	return Result{
		IgnoreParents: false,
		Screen:        screen,
		Pager:         state.PagerContext{},
	}
}

func ResultFromUIState(stt *state.UIState) Result {
	return Result{
		IgnoreParents: false,
		Screen:        nil,
		Pager:         stt.Pager,
	}
}

func EmptyResult() Result {
	return Result{
		IgnoreParents: false,
		Screen:        nil,
		Pager:         state.PagerContext{},
	}
}
