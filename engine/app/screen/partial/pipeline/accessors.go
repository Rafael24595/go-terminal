package pipeline

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
)

type StackAccessor struct {
	Get func(viewmodel.ViewModel) *stack.VStackDrawable
	Set func(viewmodel.ViewModel, *stack.VStackDrawable) viewmodel.ViewModel
}

var viewModelAccessors = map[Target]StackAccessor{
	Header: {
		Get: func(vm viewmodel.ViewModel) *stack.VStackDrawable {
			return vm.Header
		},
		Set: func(vm viewmodel.ViewModel, s *stack.VStackDrawable) viewmodel.ViewModel {
			vm.Header = s
			return vm
		},
	},
	Kernel: {
		Get: func(vm viewmodel.ViewModel) *stack.VStackDrawable {
			return vm.Kernel
		},
		Set: func(vm viewmodel.ViewModel, s *stack.VStackDrawable) viewmodel.ViewModel {
			vm.Kernel = s
			return vm
		},
	},
	Footer: {
		Get: func(vm viewmodel.ViewModel) *stack.VStackDrawable {
			return vm.Footer
		},
		Set: func(vm viewmodel.ViewModel, s *stack.VStackDrawable) viewmodel.ViewModel {
			vm.Footer = s
			return vm
		},
	},
}

func FindViewModelAccessor(target Target) (StackAccessor, bool) {
	accessor, ok := viewModelAccessors[target]
	if !ok {
		return StackAccessor{}, false
	}
	return accessor, true
}
