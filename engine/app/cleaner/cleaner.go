package cleaner

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
)

type StateCleaner struct {
	Cleanup func(screen.ScreenResult, *state.UIState) *state.UIState
}
