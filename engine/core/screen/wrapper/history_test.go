package wrapper

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/core/screen"
)

func TestHistory_ToScreen(t *testing.T) {
	base := screen.Screen{
		Name: func() string { return "base" },
		Update: func(s *state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.EmptyScreenResult()
		},
		View: func(state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state.UIState{})
		},
	}

	h := NewHistory(base)
	scrn := h.ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestHistory_BackNavigation(t *testing.T) {
	stt := &state.UIState{}

	mockBase := screen_test.MockScreen{
		Name: "base",
	}

	mockNext := screen_test.MockScreen{
		Name: "next",
		Update: func(s *state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			base := mockBase.ToScreen()
			return screen.ScreenResultFromScreen(&base)
		},
	}

	h := NewHistory(mockNext.ToScreen())
	scrn := h.ToScreen()

	assert.Equal(t, scrn.Name(), "next")

	result := scrn.Update(stt, screen.ScreenEvent{})
	assert.NotNil(t, result.Screen)
	assert.Equal(t, result.Screen.Name(), "base")

	backResult := result.Screen.Update(stt, screen.ScreenEvent{
		Key: *key.NewKeyCode(key.CustomActionBack),
	})

	assert.NotNil(t, backResult.Screen)
	assert.Equal(t, backResult.Screen.Name(), "next")
}

func TestHistory_ViewFooter(t *testing.T) {
	mock := screen_test.MockScreen{}
	scrn := mock.ToScreen()

	h := NewHistory(scrn)

	vm := h.view(*state.NewUIState())

	vm.Footer.Init(terminal.Winsize{})
	footer, _ := vm.Footer.Draw()

	assert.Equal(t, len(footer), 0)

	h.history = &scrn
	vm = h.view(*state.NewUIState())

	vm.Footer.Init(terminal.Winsize{})
	footer, _ = vm.Footer.Draw()

	assert.True(t, len(footer) > 0)
}
