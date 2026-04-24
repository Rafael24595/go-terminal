package cleaner

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

type Cleanup func(screen.ScreenResult, *state.UIState) *state.UIState

type StateCleaner struct {
	Cleanup Cleanup
}
