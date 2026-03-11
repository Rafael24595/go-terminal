package wrapper

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/test/support/assert"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/core/screen"
)

func TestHelp_ToScreen(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	h := NewHelp(mock.ToScreen())
	scrn := h.ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestHelp_ToggleHelpKey(t *testing.T) {
	called := false

	mock := screen_test.MockScreen{}

	scrn := NewHelp(mock.ToScreen()).ToScreen()

	state := &state.UIState{}
	event := screen.ScreenEvent{
		Key: *key.NewKeyCode(key.CustomActionHelp),
	}

	scrn.Update(state, event)

	assert.True(t, state.Helper.ShowHelp)
	assert.False(t, called)
}

func TestHelp_DelegatesUpdateWhenKeyRequired(t *testing.T) {
	called := false

	ky := *key.NewKeyCode(key.CustomActionHelp)

	mock := screen_test.MockScreen{
		Definition: func() screen.Definition {
			return screen.DefinitionFromKeys(ky)
		},
		Update: func(s *state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			called = true
			return screen.EmptyScreenResult()
		},
	}

	scrn := NewHelp(mock.ToScreen()).ToScreen()

	state := &state.UIState{}
	event := screen.ScreenEvent{
		Key: ky,
	}

	scrn.Update(state, event)

	assert.False(t, state.Helper.ShowHelp)
	assert.True(t, called)
}

func TestHelp_WrapsReturnedScreen(t *testing.T) {
	called := false

	ky := *key.NewKeyCode(key.ActionEnter)

	mockNext := screen_test.MockScreen{
		Name: "next",
	}

	mockBase := screen_test.MockScreen{
		Definition: func() screen.Definition {
			return screen.DefinitionFromKeys(ky)
		},
		Update: func(s *state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
			called = true
			next := mockNext.ToScreen()
			return screen.ScreenResult{
				Screen: &next,
			}
		},
	}

	help := NewHelp(mockBase.ToScreen())
	wrapped := help.ToScreen()

	stt := &state.UIState{}
	evt := screen.ScreenEvent{
		Key: ky,
	}

	result := wrapped.Update(stt, screen.ScreenEvent{
		Key: *key.NewKeyCode(key.CustomActionHelp),
	})

	assert.True(t, stt.Helper.ShowHelp)

	result = wrapped.Update(stt, evt)

	assert.True(t, called)
	assert.NotNil(t, result.Screen)
	assert.Equal(t, "next", result.Screen.Name())

	vm := result.Screen.View(state.UIState{})

	assert.True(t, vm.Helper.Show)
}
