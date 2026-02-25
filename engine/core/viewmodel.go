package core

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

var default_pager = state.NewPagePager()

type ViewModel struct {
	Header *LayerStack
	Lines  *LayerStack
	Footer *LayerStack
	Input  *InputLine
	Pager  state.PagerStrategy
}

func ViewModelFromUIState(stt state.UIState) *ViewModel {
	return &ViewModel{
		Header: NewLayerStack(),
		Lines:  NewLayerStack(),
		Footer: NewLayerStack(),
		Input:  nil,
		Pager:  default_pager,
	}
}

func (v *ViewModel) SetInput(input *InputLine) *ViewModel {
	v.Input = input
	return v
}

func (v *ViewModel) SetStrategy(strategy state.PagerStrategy) *ViewModel {
	v.Pager = strategy
	return v
}

func (s *ViewModel) IsPagerMode(mode state.PagerMode) bool {
	return s.Pager.Mode == mode
}

func (s *ViewModel) PagerMatch(state state.UIState, ctx state.PagerContext) bool {
	return s.Pager.Match(state, ctx)
}

func (v *ViewModel) InitStaticLayers(size terminal.Winsize) (*LayerStack, *LayerStack) {
	return v.Header.Init(size), v.Footer.Init(size)
}

func (v *ViewModel) InitDynamicLayers(size terminal.Winsize) *LayerStack {
	return v.Lines.Init(size)
}

func (v *ViewModel) InitInputLine(size terminal.Winsize) (*Drawable, bool) {
	if v.Input == nil {
		return nil, false
	}

	drawable := v.Input.ToDrawable()
	drawable.Init(size)

	return &drawable, true
}
