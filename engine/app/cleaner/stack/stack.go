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

func Cleanup(result screen.Result, stt *state.UIState) *state.UIState {
	if result.Node == nil {
		return stt
	}

	stt.Stack.RetainOnly(
		result.Node.Stack,
	)

	return stt
}
