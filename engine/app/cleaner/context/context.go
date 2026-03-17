package context

import (
	"github.com/Rafael24595/go-terminal/engine/app/cleaner"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
)

func NewContextCleaner() cleaner.StateCleaner {
	return cleaner.StateCleaner{
		Cleanup: cleanup,
	}
}

func cleanup(result screen.ScreenResult, stt *state.UIState) *state.UIState {
	if result.Screen != nil {
		stack := result.Screen.Stack()
		stt.Stack.RetainOnly(stack)
	}
	return stt
}
