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
	header := vm.Header.ToUnit()
	header.Drawable.Init()
	header.Drawable.Draw(winsize)

	footer := vm.Footer.ToUnit()
	footer.Drawable.Init()
	footer.Drawable.Draw(winsize)

	lines := vm.Kernel.ToUnit()
	lines.Drawable.Init()
	lines.Drawable.Draw(winsize)
}

func TestPipeline_ToNode(t *testing.T) {
	name := "base"
	mock := screen_test.MockScreen{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, name)
}

func TestPipeline_Propagate(t *testing.T) {
	name := "base"
	mock := screen_test.MockScreen{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_Propagate(t, name, 0, node)
}

func TestPipeline_WrapsReturnedScreen(t *testing.T) {
	called := false

	mockNext := screen_test.MockScreen{
		Name: "next",
	}

	mockBase := screen_test.MockScreen{
		Update: func(s *state.UIState, _ screen.Event) screen.Result {
			called = true
			next := mockNext.ToNode()
			return screen.Result{
				Node: &next,
			}
		},
	}

	help := New(mockBase.ToNode()).
		ToNode()

	stt := &state.UIState{}
	evt := screen.Event{}

	result := help.Screen.Update(stt, evt)

	assert.True(t, called)
	assert.NotNil(t, result.Node.Screen)
	assert.Equal(t, "next", result.Node.Name)
}

func TestPipeline_ActionSingleFocus(t *testing.T) {
	headerBase := drawable_test.MockUnit{
		Name: "mock_01",
	}
	kernelBase := drawable_test.MockUnit{
		Name: "mock_01",
	}
	footerBase := drawable_test.MockUnit{
		Name: "mock_01",
	}

	mockNode := screen_test.MockScreen{
		View: func(_ state.UIState) viewmodel.ViewModel {
			vm := viewmodel.NewViewModel()
			vm.Header.Push(headerBase.ToUnit())
			vm.Kernel.Push(kernelBase.ToUnit())
			vm.Footer.Push(footerBase.ToUnit())
			return *vm
		},
	}

	w := New(mockNode.ToNode())

	w.PushSteps(
		func(vm viewmodel.ViewModel) viewmodel.ViewModel {
			mock := drawable_test.MockUnit{
				Name: "mock_02",
			}

			vm.Header.Unshift(
				mock.ToUnit(),
			)
			return vm
		},
		func(vm viewmodel.ViewModel) viewmodel.ViewModel {
			mock1 := drawable_test.MockUnit{
				Name: "mock_02",
			}
			mock2 := drawable_test.MockUnit{
				Name: "mock_03",
			}

			vm.Kernel.Push(
				mock2.ToUnit(),
				mock1.ToUnit(),
			)
			return vm
		},
		func(vm viewmodel.ViewModel) viewmodel.ViewModel {
			vm.Footer = stack.NewVStack()
			return vm
		},
	)

	node := w.ToNode()
	vm := node.Screen.View(state.UIState{})

	drawSources(
		node.Screen.View(state.UIState{}),
		winsize.Winsize{},
	)

	header := vm.Header.Units()
	assert.Len(t, 2, header)
	assert.Equal(t, "mock_02", header[0].Name)

	kernel := vm.Kernel.Units()
	assert.Len(t, 3, kernel)
	assert.Equal(t, "mock_03", kernel[1].Name)
	assert.Equal(t, "mock_02", kernel[2].Name)

	assert.Len(t, 0, vm.Footer.Units())
}
