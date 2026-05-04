package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

type DefinitionFunc func() Definition
type UpdateFunc func(*state.UIState, ScreenEvent) Result
type ViewFunc func(state.UIState) viewmodel.ViewModel

type Screen struct {
	Name       string
	Stack      set.Set[string]
	Definition DefinitionFunc
	Update     UpdateFunc
	View       ViewFunc
}
