package history

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestHistory_ToScreen(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
	}

	scrn := New(base.ToScreen()).
		ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestHistory_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := New(mock.ToScreen()).
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

	scrn := New(mockNext.ToScreen()).
		ToScreen()

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

	h := New(scrn)

	vm := h.view(*state.NewUIState())

	footer := vm.Footer.ToDrawable()
	footer.Init()

	lines, _ := footer.Draw(winsize.Winsize{})

	assert.Len(t, 0, lines)

	h.history = &scrn
	vm = h.view(*state.NewUIState())

	footer = vm.Footer.ToDrawable()
	footer.Init()

	lines, _ = footer.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.Len(t, 1, lines)
}
