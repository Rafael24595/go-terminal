package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestHistory_ToScreen(t *testing.T) {
	base := screen.Screen{
		Name: func() string { return "base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.EmptyScreenResult()
		},
		View: func(state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state.UIState{})
		},
	}

	h := NewHistory(base)
	screen := h.ToScreen()

	Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestHistory_BackNavigation(t *testing.T) {
	stt := state.UIState{}

	base := screen.Screen{
		Definition: func() screen.Definition { return screen.Definition{} },
		Name:       func() string { return "base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.EmptyScreenResult()
		},
		View: func(state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state.UIState{})
		},
	}

	next := screen.Screen{
		Definition: func() screen.Definition { return screen.Definition{} },
		Name:       func() string { return "next" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.ScreenResultFromScreen(&base)
		},
		View: func(state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state.UIState{})
		},
	}

	h := NewHistory(next)
	scrn := h.ToScreen()

	assert.Equal(t, scrn.Name(), "next")

	result := scrn.Update(stt, screen.ScreenEvent{})
	assert.NotNil(t, result.Screen)
	assert.Equal(t, result.Screen.Name(), "base")

	backResult := result.Screen.Update(stt, screen.ScreenEvent{Key: *key.NewKeyRune('b')})
	assert.NotNil(t, backResult.Screen)
	assert.Equal(t, backResult.Screen.Name(), "next")
}

func TestHistory_ViewFooter(t *testing.T) {
	base := screen.Screen{
		Name: func() string {
			return "base"
		},
		View: func(state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state.UIState{})
		},
	}

	h := NewHistory(base)

	vm := h.view(*state.NewUIState())

	vm.Footer.Init(terminal.Winsize{})
	footer, _ := vm.Footer.Draw()

	assert.Equal(t, len(footer), 0)

	h.history = &base
	vm = h.view(*state.NewUIState())
	footer, _ = vm.Footer.Draw()
	
	assert.True(t, len(footer) > 0)
}
