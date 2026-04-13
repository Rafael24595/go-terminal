package viewmodel

import (
	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/model/help"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	drawable_help "github.com/Rafael24595/go-terminal/engine/layout/drawable/help"
)

type ViewModel struct {
	Header *stack.VStackDrawable
	Kernel *stack.VStackDrawable
	Footer *stack.VStackDrawable
	Input  *InputLine
	Pager  pager.PagerStrategy
	Helper *help.HelpMeta
}

func NewViewModel() *ViewModel {
	return &ViewModel{
		Header: stack.NewVStackDrawable(),
		Kernel: stack.NewVStackDrawable(),
		Footer: stack.NewVStackDrawable(),
		Input:  nil,
		Pager:  pager.NewStrategy(),
		Helper: help.NewHelpMeta(),
	}
}

func (v *ViewModel) SetInput(input *InputLine) *ViewModel {
	v.Input = input
	return v
}

func (v *ViewModel) InitStaticLayers() (drawable.Drawable, drawable.Drawable) {
	header := v.Header.ToDrawable()
	header.Init()

	footer := v.Footer.ToDrawable()
	footer.Init()

	return header, footer
}

func (v *ViewModel) InitDynamicLayers(size terminal.Winsize) drawable.Drawable {
	kernel := v.Kernel.ToDrawable()
	kernel.Init()

	return kernel
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

func (v *ViewModel) Clone() *ViewModel {
	vm := NewViewModel()

	vm.Header.Push(v.Header.Items()...)
	vm.Kernel.Push(v.Kernel.Items()...)
	vm.Footer.Push(v.Footer.Items()...)
	vm.Input = v.Input
	vm.Pager = v.Pager
	vm.Helper.Push(v.Helper.Fields...)

	return vm
}
