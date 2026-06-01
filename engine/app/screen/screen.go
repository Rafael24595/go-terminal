package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

const (
	SystemMetaTag = "system_meta"
)

type KeysFunc func() Definition
type TickFunc func(*state.UIState, Event) Result
type ViewFunc func(state.UIState) viewmodel.ViewModel

type Screen struct {
	Keys KeysFunc
	Tick TickFunc
	View ViewFunc
}

func IsZeroScreen(screen Screen) bool {
	if screen.Keys == nil {
		return true
	}

	if screen.Tick == nil {
		return true
	}

	if screen.View == nil {
		return true
	}

	return false
}
