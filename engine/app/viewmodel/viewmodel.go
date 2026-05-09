package viewmodel

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// TODO: Use Screen and Drawable sources to manage Header and Footer.
type ViewModel struct {
	Header   *stack.VStackDrawable
	Kernel   *stack.VStackDrawable
	Footer   *stack.VStackDrawable
	Pager    *pager.PagerStrategy
	Behavior BehaviorContext
}

func NewViewModel() *ViewModel {
	return &ViewModel{
		Header:   stack.NewVStack(),
		Kernel:   stack.NewVStack(),
		Footer:   stack.NewVStack(),
		Pager:    pager.NewStrategy(),
		Behavior: BehaviorContext{},
	}
}

func (v *ViewModel) InitStaticLayers() (drawable.Drawable, drawable.Drawable) {
	header := v.Header.ToDrawable()
	header.Init()

	footer := v.Footer.ToDrawable()
	footer.Init()

	return header, footer
}

func (v *ViewModel) InitDynamicLayers(size winsize.Winsize) drawable.Drawable {
	kernel := v.Kernel.ToDrawable()
	kernel.Init()

	return kernel
}

func (v *ViewModel) Clone() *ViewModel {
	vm := NewViewModel()

	vm.Header.Push(v.Header.Items()...)
	vm.Kernel.Push(v.Kernel.Items()...)
	vm.Footer.Push(v.Footer.Items()...)
	vm.Pager = v.Pager

	return vm
}
