package wrapper

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestPagination_ToScreen(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
	}

	p := NewPagination(base.ToScreen())
	screen := p.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestPagination_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := NewPagination(mock.ToScreen()).
		ToScreen().
		Stack()

	assert.True(t, stack.Has("base"))
}

func TestPagination_LocalUpdate(t *testing.T) {
	stt := state.NewUIState()
	base := screen_test.MockScreen{
		Name: "base",
	}

	p := NewPagination(base.ToScreen())

	scrn := p.ToScreen()

	stt.Pager.TargetPage = 0
	result := scrn.Update(stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowLeft)})
	assert.Equal(t, result.Pager.TargetPage, 0)

	result = scrn.Update(stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowRight)})
	assert.Equal(t, result.Pager.TargetPage, 1)
}

func TestPagination_ViewFooter(t *testing.T) {
	stt := state.NewUIState()
	stt.Pager.ActualPage = 3

	base := screen_test.MockScreen{
		Name: "base",
		View: func(_ state.UIState) viewmodel.ViewModel {
			vm := viewmodel.NewViewModel()
			vm.Pager.SetPredicate(pager.PredicatePage())
			return *vm
		},
	}

	p := NewPagination(base.ToScreen())
	vm := p.view(*stt)

	footer := vm.Footer.ToDrawable()
	footer.Init()

	lines, _ := footer.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.True(t, len(lines) > 0)
	assert.Contains(t, text.LineToString(&lines[0]), "page: 3")
}

func TestPagination_UpdateDelegates(t *testing.T) {
	called := false

	base := screen_test.MockScreen{
		Name: "base",
		Update: func(s *state.UIState, e screen.ScreenEvent) screen.ScreenResult {
			called = true
			return screen.EmptyScreenResult()
		},
	}

	p := NewPagination(base.ToScreen())
	scrn := p.ToScreen()

	scrn.Update(state.NewUIState(), screen.ScreenEvent{Key: *key.NewKeyRune('x')})

	assert.True(t, called, "screen.Update should be called")
}

func TestPagination_PageNeverNegative(t *testing.T) {
	stt := state.NewUIState()
	stt.Pager.TargetPage = 0

	base := screen_test.MockScreen{
		Name: "base",
	}

	p := NewPagination(base.ToScreen())
	scrn := p.ToScreen()

	scrn.Update(stt, screen.ScreenEvent{Key: *key.NewKeyCode(key.ActionArrowLeft)})
	assert.Equal(t, stt.Pager.TargetPage, 0)
}
