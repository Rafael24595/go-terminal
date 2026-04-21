package cleaner

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
)

type Cleanup func(screen.ScreenResult, *state.UIState) *state.UIState

type StateCleaner struct {
	Cleanup Cleanup
}
