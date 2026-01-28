package wrapper_commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestHistory_ToScreen(t *testing.T) {
	base := core.Screen{
		Name: func() string { return "base" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult {
			return core.ScreenResultFromState(s)
		},
		View: func(state.UIState) core.ViewModel {
			return core.ViewModel{}
		},
	}

	h := NewHistory(base)
	screen := h.ToScreen()

	Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestHistory_BackNavigation(t *testing.T) {
	stt := state.UIState{}

	base := core.Screen{
		Name: func() string { return "base" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult {
			return core.ScreenResultFromState(s)
		},
		View: func(state.UIState) core.ViewModel {
			return core.ViewModel{}
		},
	}

	next := core.Screen{
		Name: func() string { return "next" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult {
			return core.NewScreenResult(s, &base)
		},
		View: func(state.UIState) core.ViewModel {
			return core.ViewModel{}
		},
	}

	h := NewHistory(next)
	screen := h.ToScreen()

	assert.Equal(t, screen.Name(), "next")

	result := screen.Update(stt, core.ScreenEvent{})
	assert.NotNil(t, result.Screen)
	assert.Equal(t, result.Screen.Name(), "base")

	backResult := result.Screen.Update(stt, core.ScreenEvent{Key: "b"})
	assert.NotNil(t, backResult.Screen)
	assert.Equal(t, backResult.Screen.Name(), "next")
}

func TestHistory_ViewFooter(t *testing.T) {
	base := core.Screen{
		Name: func() string { return "base" },
		View: func(state.UIState) core.ViewModel {
			return core.ViewModel{}
		},
	}

	h := NewHistory(base)

	vm := h.view(*state.NewUIState())
	assert.Equal(t, len(vm.Footer), 0)

	h.history = &base
	vm = h.view(*state.NewUIState())
	assert.True(t, len(vm.Footer) > 0)
}
