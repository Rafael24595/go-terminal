package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

const (
	SystemMetaTag = "system_meta"
)

type DefinitionFunc func() Definition
type UpdateFunc func(*state.UIState, Event) Result
type ViewFunc func(state.UIState) viewmodel.ViewModel

type Screen struct {
	Definition DefinitionFunc
	Update     UpdateFunc
	View       ViewFunc
}

func IsZeroScreen(screen Screen) bool {
	if screen.Definition == nil {
		return true
	}

	if screen.Update == nil {
		return true
	}

	if screen.View == nil {
		return true
	}

	return false
}
