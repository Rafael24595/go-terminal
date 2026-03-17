package partial

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/render/spacer"
	"github.com/Rafael24595/go-terminal/test/support/assert"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestSpacer_ToScreen(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
	}

	h := NewSpacer(base.ToScreen())
	scrn := h.ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestSpacer_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := NewSpacer(mock.ToScreen()).
		ToScreen().
		Stack()

	assert.True(t, stack.Has("base"))
}

func TestSpacer_WrapsNextScreen(t *testing.T) {
	next := screen_test.MockScreen{
		Name: "next",
	}

	base := screen_test.MockScreen{
		Update: func(*state.UIState, screen.ScreenEvent) screen.ScreenResult {
			scr := next.ToScreen()
			return screen.ScreenResult{
				Screen: &scr,
			}
		},
	}

	w := NewSpacer(base.ToScreen())

	s := w.ToScreen()

	result := s.Update(&state.UIState{}, screen.ScreenEvent{})

	assert.NotNil(t, result.Screen)
	assert.Equal(t, "next", result.Screen.Name())
}

func TestSpacer_AddsHeaderLinesWhenEmpty(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
	}

	w := NewSpacer(base.ToScreen()).
		Header(spacer.NewSpacerMeta(1, spacer.SpacerAppend))

	s := w.ToScreen()

	vm := s.View(state.UIState{})

	items := vm.Header.Items()

	assert.Len(t, 0, items)
}

func TestSpacer_AddsHeaderLines(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
		View: func(stt state.UIState) viewmodel.ViewModel {
			mc := drawable_test.MockDrawable{}
			dw := mc.ToDrawable()
			dw.Name = "mock_header"
			vm := viewmodel.ViewModelFromUIState(stt)
			vm.Header.Shift(dw)
			return *vm
		},
	}

	w := NewSpacer(base.ToScreen()).
		Header(spacer.NewSpacerMeta(1, spacer.SpacerAppend))

	s := w.ToScreen()

	vm := s.View(state.UIState{})

	items := vm.Header.Items()

	assert.Len(t, 2, items)

	assert.Equal(t, "mock_header", items[0].Name)
	assert.NotEqual(t, "mock_header", items[1].Name)
}

func TestSpacer_HeaderBetween(t *testing.T) {
	base := screen_test.MockScreen{
		View: func(stt state.UIState) viewmodel.ViewModel {
			vm := viewmodel.ViewModelFromUIState(stt)

			m1 := drawable_test.MockDrawable{}
			d1 := m1.ToDrawable()
			d1.Name = "h1"

			m2 := drawable_test.MockDrawable{}
			d2 := m2.ToDrawable()
			d2.Name = "h2"

			vm.Header.Shift(d1, d2)
			return *vm
		},
	}

	w := NewSpacer(base.ToScreen()).
		Header(spacer.NewSpacerMeta(1, spacer.SpacerAfterEach))

	s := w.ToScreen()

	vm := s.View(state.UIState{})

	items := vm.Header.Items()

	assert.Len(t, 4, items)

	assert.Equal(t, "h1", items[0].Name)
	assert.NotEqual(t, "h1", items[1].Name)

	assert.Equal(t, "h2", items[2].Name)
	assert.NotEqual(t, "h2", items[3].Name)
}

func TestSpacer_AddsFooterLinesWhenEmpty(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
	}

	w := NewSpacer(base.ToScreen()).
		Footer(spacer.NewSpacerMeta(1, spacer.SpacerAppend))

	s := w.ToScreen()

	vm := s.View(state.UIState{})

	items := vm.Header.Items()

	assert.Len(t, 0, items)
}

func TestSpacer_AddsFooterLines(t *testing.T) {
	base := screen_test.MockScreen{
		Name: "base",
		View: func(stt state.UIState) viewmodel.ViewModel {
			mc := drawable_test.MockDrawable{}
			dw := mc.ToDrawable()
			dw.Name = "mock_footer"
			vm := viewmodel.ViewModelFromUIState(stt)
			vm.Footer.Shift(dw)
			return *vm
		},
	}

	w := NewSpacer(base.ToScreen()).
		Footer(spacer.NewSpacerMeta(1, spacer.SpacerAppend))

	s := w.ToScreen()

	vm := s.View(state.UIState{})

	items := vm.Footer.Items()

	assert.Len(t, 2, items)

	assert.NotEqual(t, "mock_footer", items[0].Name)
	assert.Equal(t, "mock_footer", items[1].Name)
}

func TestSpacer_FooterBetween(t *testing.T) {
	base := screen_test.MockScreen{
		View: func(stt state.UIState) viewmodel.ViewModel {
			vm := viewmodel.ViewModelFromUIState(stt)

			m1 := drawable_test.MockDrawable{}
			d1 := m1.ToDrawable()
			d1.Name = "f1"

			m2 := drawable_test.MockDrawable{}
			d2 := m2.ToDrawable()
			d2.Name = "f2"

			vm.Footer.Shift(d1, d2)
			return *vm
		},
	}

	w := NewSpacer(base.ToScreen()).
		Footer(spacer.NewSpacerMeta(1, spacer.SpacerAfterEach))

	s := w.ToScreen()

	vm := s.View(state.UIState{})

	items := vm.Footer.Items()

	assert.Len(t, 4, items)

	assert.NotEqual(t, "f1", items[0].Name)
	assert.Equal(t, "f1", items[1].Name)

	assert.NotEqual(t, "f2", items[2].Name)
	assert.Equal(t, "f2", items[3].Name)
}
