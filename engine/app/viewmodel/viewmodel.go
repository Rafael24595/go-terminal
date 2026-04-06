package viewmodel

import (
	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/model/help"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	drawable_help "github.com/Rafael24595/go-terminal/engine/layout/drawable/help"
)

type ViewModel struct {
	Header *stack.StackDrawable
	Kernel *stack.StackDrawable
	Footer *stack.StackDrawable
	Input  *InputLine
	Pager  pager.PagerStrategy
	Helper *help.HelpMeta
}

func ViewModelFromUIState(stt state.UIState) *ViewModel {
	return &ViewModel{
		Header: stack.NewStackDrawable(),
		Kernel: stack.NewStackDrawable(),
		Footer: stack.NewStackDrawable(),
		Input:  nil,
		Pager:  pager.NewStrategy(),
		Helper: help.NewHelpMeta(),
	}
}

func (v *ViewModel) SetInput(input *InputLine) *ViewModel {
	v.Input = input
	return v
}

func (v *ViewModel) InitStaticLayers() (*stack.StackDrawable, *stack.StackDrawable) {
	return v.Header.Init(), v.Footer.Init()
}

func (v *ViewModel) InitDynamicLayers(size terminal.Winsize) *stack.StackDrawable {
	return v.Kernel.Init()
}

func (v *ViewModel) InitInputLine(size terminal.Winsize) (drawable.Drawable, bool) {
	if v.Input == nil {
		return drawable.Drawable{}, false
	}

	drawable := v.Input.ToDrawable()
	drawable.Init()

	return drawable, true
}

func (v *ViewModel) InitHelper(size terminal.Winsize) (drawable.Drawable, bool) {
	if !v.Helper.Show {
		return drawable.Drawable{}, false
	}

	drawable := drawable_help.HelpDrawableFromMeta(v.Helper)
	drawable.Init()

	return drawable, true
}
