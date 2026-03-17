package viewmodel

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/model/help"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	drawable_help "github.com/Rafael24595/go-terminal/engine/layout/drawable/help"
)

var default_pager = state.NewPagePager()

type ViewModel struct {
	Header *stack.StackDrawable
	Lines  *stack.StackDrawable
	Footer *stack.StackDrawable
	Input  *InputLine
	Pager  state.PagerStrategy
	Helper *help.HelpMeta
}

func ViewModelFromUIState(stt state.UIState) *ViewModel {
	return &ViewModel{
		Header: stack.NewStackDrawable(),
		Lines:  stack.NewStackDrawable(),
		Footer: stack.NewStackDrawable(),
		Input:  nil,
		Pager:  default_pager,
		Helper: help.NewHelpMeta(),
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

func (v *ViewModel) InitStaticLayers(size terminal.Winsize) (*stack.StackDrawable, *stack.StackDrawable) {
	return v.Header.Init(size), v.Footer.Init(size)
}

func (v *ViewModel) InitDynamicLayers(size terminal.Winsize) *stack.StackDrawable {
	return v.Lines.Init(size)
}

func (v *ViewModel) InitInputLine(size terminal.Winsize) (drawable.Drawable, bool) {
	if v.Input == nil {
		return drawable.Drawable{}, false
	}

	drawable := v.Input.ToDrawable()
	drawable.Init(size)

	return drawable, true
}

func (v *ViewModel) InitHelper(size terminal.Winsize) (drawable.Drawable, bool) {
	if !v.Helper.Show {
		return drawable.Drawable{}, false
	}

	drawable := drawable_help.HelpDrawableFromMeta(v.Helper)
	drawable.Init(size)

	return drawable, true
}
