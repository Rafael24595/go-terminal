package screen_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

type MockScreen struct {
	Name       string
	Definition *screen.Definition
	Update     func(*state.UIState, screen.ScreenEvent) screen.Result
	View       func(state.UIState) viewmodel.ViewModel
	Stack      set.Set[string]
}

func (t MockScreen) ToNode() screen.Node {
	stack := t.Stack
	if t.Stack == nil {
		stack = set.SetFrom(t.Name)
	}

	return screen.NewBuilder().
		Name(t.Name).
		AddStack(stack).
		Definition(
			func() screen.Definition {
				if t.Definition != nil {
					return *t.Definition
				}

				return screen.DefinitionFromKeys()
			}).
		Update(
			func(s *state.UIState, e screen.ScreenEvent) screen.Result {
				if t.Update != nil {
					return t.Update(s, e)
				}

				return screen.ResultFromUIState(s)
			},
		).
		View(
			func(s state.UIState) viewmodel.ViewModel {
				if t.View != nil {
					return t.View(s)
				}

				return *viewmodel.NewViewModel()
			},
		).
		ToNode()
}

func Helper_ToNode(t *testing.T, node screen.Node) {
	t.Helper()

	assert.NotNil(t, node.Screen.Name, "Screen.Name")
	assert.NotNil(t, node.Screen.View, "Screen.View should be set")
	assert.NotNil(t, node.Screen.Update, "Screen.Update should be set")

	assert.NotNil(t, node.Stack, "Screen.Stack should be set")
}

func Helper_Propagate(t *testing.T, name string, child uint, node screen.Node) {
	assert.GreaterOrEqual(t, child+1, uint(len(node.Children())))
	assert.True(t, node.Stack.Has(name))
	assert.Equal(t, name, node.Children()[child].Screen.Name)
}
