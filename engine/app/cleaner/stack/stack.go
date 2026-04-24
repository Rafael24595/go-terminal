package stack

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func NewCleaner() cleaner.StateCleaner {
	return cleaner.StateCleaner{
		Cleanup: Cleanup,
	}
}

func Cleanup(res screen.ScreenResult, stt *state.UIState) *state.UIState {
	if res.Screen != nil {
		stack := res.Screen.Stack()
		stt.Stack.RetainOnly(stack)
	}
	return stt
}
