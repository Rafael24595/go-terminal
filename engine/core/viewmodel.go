package core

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type InputLine struct {
	Prompt string
	Value  string
	Cursor int
}

type ViewModel struct {
	Header *LayerStack
	Lines  *LayerStack
	Footer *LayerStack
	Input  *InputLine
	Pager  state.PagerState
	Cursor state.CursorState
}

func ViewModelFromUIState(state state.UIState) *ViewModel {
	return &ViewModel{
		Header: NewLayerStack(),
		Lines:  NewLayerStack(),
		Footer: NewLayerStack(),
		Pager:  state.Pager,
		Cursor: state.Cursor,
	}
}

func (v *ViewModel) SetPager(pager state.PagerState) *ViewModel {
	v.Pager = pager
	return v
}

func (v *ViewModel) SetCursor(cursor state.CursorState) *ViewModel {
	v.Cursor = cursor
	return v
}

func (v *ViewModel) InitStaticLayers(size terminal.Winsize) (*LayerStack, *LayerStack) {
	return v.Header.Init(size), v.Footer.Init(size)
}

func (v *ViewModel) InitDynamicLayers(size terminal.Winsize) *LayerStack {
	return v.Lines.Init(size)
}
