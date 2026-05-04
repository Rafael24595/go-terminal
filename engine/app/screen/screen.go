package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

type DefinitionFunc func() Definition
type UpdateFunc func(*state.UIState, ScreenEvent) Result
type ViewFunc func(state.UIState) viewmodel.ViewModel

type Screen struct {
	Name       string
	Definition DefinitionFunc
	Update     UpdateFunc
	View       ViewFunc
}
