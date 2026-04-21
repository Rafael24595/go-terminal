package composite

import (
	"github.com/Rafael24595/go-terminal/engine/app/cleaner"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
)

func NewCleaner(cls ...cleaner.Cleanup) cleaner.StateCleaner {
	return cleaner.StateCleaner{
		Cleanup: cleanup(cls),
	}
}

func cleanup(cls []cleaner.Cleanup) cleaner.Cleanup {
	return func(res screen.ScreenResult, stt *state.UIState) *state.UIState {
		for _, part := range cls {
			stt = part(res, stt)
		}
		return stt
	}
}
