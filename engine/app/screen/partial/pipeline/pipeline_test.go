package pipeline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func drawSources(vm viewmodel.ViewModel, winsize winsize.Winsize) {
	header := vm.Header.ToDrawable()
	header.Init()
	header.Draw(winsize)

	footer := vm.Footer.ToDrawable()
	footer.Init()
	footer.Draw(winsize)

	lines := vm.Kernel.ToDrawable()
	lines.Init()
	lines.Draw(winsize)
}

func TestPipeline_ToScreen(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	scrn := New(mock.ToScreen()).
		ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestPipeline_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := New(mock.ToScreen()).
		ToScreen().
		Stack()

	assert.True(t, stack.Has("base"))
}

func TestPipeline_WrapsReturnedScreen(t *testing.T) {
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

	help := New(mockBase.ToScreen()).
		ToScreen()

	stt := &state.UIState{}
	evt := screen.ScreenEvent{}

	result := help.Update(stt, evt)

	assert.True(t, called)
	assert.NotNil(t, result.Screen)
	assert.Equal(t, "next", result.Screen.Name())
}

func TestPipeline_ActionSingleFocus(t *testing.T) {
	headerBase := drawable_test.MockDrawable{
		Code: "mock_01",
	}
	kernelBase := drawable_test.MockDrawable{
		Code: "mock_01",
	}
	footerBase := drawable_test.MockDrawable{
		Code: "mock_01",
	}

	mockScreen := screen_test.MockScreen{
		View: func(_ state.UIState) viewmodel.ViewModel {
			vm := viewmodel.NewViewModel()
			vm.Header.Push(headerBase.ToDrawable())
			vm.Kernel.Push(kernelBase.ToDrawable())
			vm.Footer.Push(footerBase.ToDrawable())
			return *vm
		},
	}

	w := New(mockScreen.ToScreen())

	w.PushSteps(
		func(vm viewmodel.ViewModel) viewmodel.ViewModel {
			mock := drawable_test.MockDrawable{Code: "mock_02"}
			vm.Header.Unshift(
				mock.ToDrawable(),
			)
			return vm
		},
		func(vm viewmodel.ViewModel) viewmodel.ViewModel {
			mock1 := drawable_test.MockDrawable{Code: "mock_02"}
			mock2 := drawable_test.MockDrawable{Code: "mock_03"}
			vm.Kernel.Push(
				mock2.ToDrawable(),
				mock1.ToDrawable(),
			)
			return vm
		},
		func(vm viewmodel.ViewModel) viewmodel.ViewModel {
			vm.Footer = stack.NewVStackDrawable()
			return vm
		},
	)

	scn := w.ToScreen()
	vm := scn.View(state.UIState{})

	drawSources(
		scn.View(state.UIState{}),
		winsize.Winsize{},
	)

	header := vm.Header.Items()
	assert.Len(t, 2, header)
	assert.Equal(t, "mock_02", header[0].Code)

	kernel := vm.Kernel.Items()
	assert.Len(t, 3, kernel)
	assert.Equal(t, "mock_03", kernel[1].Code)
	assert.Equal(t, "mock_02", kernel[2].Code)

	assert.Len(t, 0, vm.Footer.Items())
}
