package wrapper_commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/test/support/assert"
	wrapper_terminal "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

func TestPagination_ToScreen(t *testing.T) {
	base := core.Screen{
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult {
			return core.ScreenResultFromState(s)
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
	base := core.Screen{
		Name:   func() string { return "Base" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult { return core.ScreenResultFromState(s) },
		View:   func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)

	screen := p.ToScreen()

	stt.Layout.Page = 0
	result := screen.Update(*stt, core.ScreenEvent{Key: wrapper_terminal.ARROW_LEFT})
	assert.Equal(t, result.State.Layout.Page, 0)

	result = screen.Update(*stt, core.ScreenEvent{Key: wrapper_terminal.ARROW_RIGHT})
	assert.Equal(t, result.State.Layout.Page, 1)
}

func TestPagination_ViewFooter(t *testing.T) {
	stt := state.NewUIState()
	stt.Layout.Pagination = true
	stt.Layout.Page = 3

	base := core.Screen{
		Name:   func() string { return "Base" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult { return core.ScreenResultFromState(s) },
		View:   func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)
	vm := p.View(*stt)

	assert.True(t, len(vm.Footer) > 0)
	assert.Contains(t, vm.Footer[1].Text[0].Text, "page: 3")
}

func TestPagination_UpdateDelegates(t *testing.T) {
	called := false

	base := core.Screen{
		Name: func() string { return "Base" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult {
			called = true
			return core.ScreenResultFromState(s)
		},
		View: func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)
	screen := p.ToScreen()

	screen.Update(*state.NewUIState(), core.ScreenEvent{Key: "x"})

	assert.True(t, called, "screen.Update should be called")
}

func TestPagination_PageNeverNegative(t *testing.T) {
	stt := state.NewUIState()
	stt.Layout.Page = 0

	base := core.Screen{
		Name:   func() string { return "Base" },
		Update: func(s state.UIState, e core.ScreenEvent) core.ScreenResult { return core.ScreenResultFromState(s) },
		View:   func(state state.UIState) core.ViewModel { return core.ViewModel{} },
	}

	p := NewPagination(base)
	screen := p.ToScreen()

	screen.Update(*stt, core.ScreenEvent{Key: wrapper_terminal.ARROW_LEFT})
	assert.Equal(t, stt.Layout.Page, 0)
}
