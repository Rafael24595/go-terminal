package primitive

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/model/action"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func drawSources(vm viewmodel.ViewModel, winsize terminal.Winsize) {
	header := vm.Header.ToDrawable()
	header.Init(winsize)
	header.Draw()

	footer := vm.Footer.ToDrawable()
	footer.Init(winsize)
	footer.Draw()

	lines := vm.Kernel.ToDrawable()
	lines.Init(winsize)
	lines.Draw()
}

func TestMapScreen_ToScreen(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	h := NewMapScreen(mock.ToScreen())
	scrn := h.ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestMapScreen_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := NewMapScreen(mock.ToScreen()).
		ToScreen().
		Stack()

	assert.True(t, stack.Has("base"))
}

func TestMapScreen_WrapsReturnedScreen(t *testing.T) {
	called := false

	mockNext := screen_test.MockScreen{
		Name: "next",
	}

	mockBase := screen_test.MockScreen{
		Update: func(s *state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
			called = true
			next := mockNext.ToScreen()
			return screen.ScreenResult{
				Screen: &next,
			}
		},
	}

	help := NewMapScreen(mockBase.ToScreen()).
		ToScreen()

	stt := &state.UIState{}
	evt := screen.ScreenEvent{}

	result := help.Update(stt, evt)

	assert.True(t, called)
	assert.NotNil(t, result.Screen)
	assert.Equal(t, "next", result.Screen.Name())
}

func TestMapScreen_ActionSingleFocus(t *testing.T) {
	headerBase := drawable_test.MockDrawable{}
	linesBase := drawable_test.MockDrawable{}
	footerBase := drawable_test.MockDrawable{}

	headerMock := drawable_test.MockDrawable{}

	mockScreen := screen_test.MockScreen{
		View: func(stt state.UIState) viewmodel.ViewModel {
			vm := viewmodel.ViewModelFromUIState(stt)
			vm.Header.Push(headerBase.ToDrawable())
			vm.Kernel.Push(linesBase.ToDrawable())
			vm.Footer.Push(footerBase.ToDrawable())
			return *vm
		},
	}

	w := NewMapScreen(mockScreen.ToScreen())

	w.PushAction(action.NewAction(
		action.ActionMapGroup,
		action.FocusHeader,
		func(d ...drawable.Drawable) []drawable.Drawable {
			return []drawable.Drawable{
				headerMock.ToDrawable(),
			}
		},
	))

	s := w.ToScreen()

	drawSources(
		s.View(state.UIState{}),
		terminal.Winsize{},
	)

	assert.False(t, headerBase.InitCalled)

	assert.True(t, headerMock.InitCalled)
	assert.True(t, linesBase.InitCalled)
	assert.True(t, footerBase.InitCalled)
}

func TestMapScreen_ActionMultipleFocus(t *testing.T) {
	headerBase := drawable_test.MockDrawable{}
	linesBase := drawable_test.MockDrawable{}
	footerBase := drawable_test.MockDrawable{}

	multipleMock := drawable_test.MockDrawable{}

	mockScreen := screen_test.MockScreen{
		View: func(stt state.UIState) viewmodel.ViewModel {
			vm := viewmodel.ViewModelFromUIState(stt)
			vm.Header.Push(headerBase.ToDrawable())
			vm.Kernel.Push(linesBase.ToDrawable())
			vm.Footer.Push(footerBase.ToDrawable())
			return *vm
		},
	}

	w := NewMapScreen(mockScreen.ToScreen())

	w.PushAction(action.NewAction(
		action.ActionMapGroup,
		action.MergeFocus(action.FocusBody, action.FocusFooter),
		func(d ...drawable.Drawable) []drawable.Drawable {
			return []drawable.Drawable{
				multipleMock.ToDrawable(),
			}
		},
	))

	s := w.ToScreen()

	drawSources(
		s.View(state.UIState{}),
		terminal.Winsize{},
	)

	assert.False(t, linesBase.InitCalled)
	assert.False(t, footerBase.InitCalled)

	assert.Equal(t, 2, multipleMock.DrawCalls)
}
