package viewmodel

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
)

// TODO: Use Screen and Units sources to manage Header and Footer.
type ViewModel struct {
	Header   *stack.VStackUnit
	Kernel   *stack.VStackUnit
	Footer   *stack.VStackUnit
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

func (v *ViewModel) InitStaticLayers() (drawable.Unit, drawable.Unit) {
	header := v.Header.ToUnit()
	header.Drawable.Init()

	footer := v.Footer.ToUnit()
	footer.Drawable.Init()

	return header, footer
}

func (v *ViewModel) InitDynamicLayers() drawable.Unit {
	kernel := v.Kernel.ToUnit()
	kernel.Drawable.Init()

	return kernel
}

func (v *ViewModel) Clone() *ViewModel {
	vm := NewViewModel()

	vm.Header.Push(v.Header.Units()...)
	vm.Kernel.Push(v.Kernel.Units()...)
	vm.Footer.Push(v.Footer.Units()...)
	vm.Pager = v.Pager

	return vm
}
