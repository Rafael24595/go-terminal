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

func TestPagination_ToScreen(t *testing.T) {
	base := screen.Screen{
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.EmptyScreenResult()
		},
		View: func(state state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state)
		},
	}

	p := NewPagination(base)
	screen := p.ToScreen()

	Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "Base")
}

func TestPagination_LocalUpdate(t *testing.T) {
	stt := state.NewUIState()
	base := screen.Screen{
		Definition: func() screen.Definition { return screen.Definition{} },
		Name:       func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.EmptyScreenResult()
		},
		View: func(state state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state)
		},
	}

	p := NewPagination(base)

	scrn := p.ToScreen()

	stt.Pager.Page = 0
	result := scrn.Update(*stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowLeft)})
	assert.Equal(t, result.Pager.Page, 0)

	result = scrn.Update(*stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowRight)})
	assert.Equal(t, result.Pager.Page, 1)
}

func TestPagination_ViewFooter(t *testing.T) {
	stt := state.NewUIState()
	stt.Pager.Enabled = true
	stt.Pager.Page = 3

	base := screen.Screen{
		Definition: func() screen.Definition { return screen.Definition{} },
		Name:       func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.EmptyScreenResult()
		},
		View: func(state state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state)
		},
	}

	p := NewPagination(base)
	vm := p.View(*stt)

	vm.Footer.Init(terminal.Winsize{ Cols: 10 })
	footer, _ := vm.Footer.Draw()

	assert.True(t, len(footer) > 0)
	assert.Contains(t, footer[1].String(), "page: 3")
}

func TestPagination_UpdateDelegates(t *testing.T) {
	called := false

	base := screen.Screen{
		Definition: func() screen.Definition { return screen.Definition{} },
		Name:       func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			called = true
			return screen.EmptyScreenResult()
		},
		View: func(state state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state)
		},
	}

	p := NewPagination(base)
	scrn := p.ToScreen()

	scrn.Update(*state.NewUIState(), screen.ScreenEvent{Key: *key.NewKeyRune('x')})

	assert.True(t, called, "screen.Update should be called")
}

func TestPagination_PageNeverNegative(t *testing.T) {
	stt := state.NewUIState()
	stt.Pager.Page = 0

	base := screen.Screen{
		Definition: func() screen.Definition { return screen.Definition{} },
		Name:       func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.EmptyScreenResult()
		},
		View: func(state state.UIState) core.ViewModel {
			return *core.ViewModelFromUIState(state)
		},
	}

	p := NewPagination(base)
	scrn := p.ToScreen()

	scrn.Update(*stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowLeft)})
	assert.Equal(t, stt.Pager.Page, 0)
}
