package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestPagination_ToScreen(t *testing.T) {
	base := screen.Screen{
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.ScreenResultFromState(s)
		},
		View: func(state state.UIState) core.ViewModel {
			return core.ViewModel{}
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
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.ScreenResultFromState(s)
		},
		View: func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)

	scrn := p.ToScreen()

	stt.Layout.Page = 0
	result := scrn.Update(*stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.KeyArrowLeft)})
	assert.Equal(t, result.State.Layout.Page, 0)

	result = scrn.Update(*stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.KeyArrowRight)})
	assert.Equal(t, result.State.Layout.Page, 1)
}

func TestPagination_ViewFooter(t *testing.T) {
	stt := state.NewUIState()
	stt.Layout.Pagination = true
	stt.Layout.Page = 3

	base := screen.Screen{
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.ScreenResultFromState(s)
		},
		View: func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)
	vm := p.View(*stt)

	assert.True(t, len(vm.Footer) > 0)
	assert.Contains(t, vm.Footer[1].Text[0].Text, "page: 3")
}

func TestPagination_UpdateDelegates(t *testing.T) {
	called := false

	base := screen.Screen{
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			called = true
			return screen.ScreenResultFromState(s)
		},
		View: func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)
	scrn := p.ToScreen()

	scrn.Update(*state.NewUIState(), screen.ScreenEvent{Key: *key.NewKeyRune('x')})

	assert.True(t, called, "screen.Update should be called")
}

func TestPagination_PageNeverNegative(t *testing.T) {
	stt := state.NewUIState()
	stt.Layout.Page = 0

	base := screen.Screen{
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			return screen.ScreenResultFromState(s)
		},
		View: func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)
	scrn := p.ToScreen()

	scrn.Update(*stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.KeyArrowLeft)})
	assert.Equal(t, stt.Layout.Page, 0)
}
