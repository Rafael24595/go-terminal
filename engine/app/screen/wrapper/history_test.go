package wrapper

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestHistory_ToScreen(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
	}

	h := NewHistory(base.ToScreen())
	scrn := h.ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestHistory_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := NewHistory(mock.ToScreen()).
		ToScreen().
		Stack()

	assert.True(t, stack.Has("base"))
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

	vm.Footer.Init()
	footer, _ := vm.Footer.Draw(terminal.Winsize{})

	assert.Len(t, 0, footer)

	h.history = &scrn
	vm = h.view(*state.NewUIState())

	vm.Footer.Init()
	footer, _ = vm.Footer.Draw(terminal.Winsize{
		Rows: 3,
	})

	assert.True(t, len(footer) > 0)
}
