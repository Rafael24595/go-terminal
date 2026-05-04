package screen

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

func TestBuilder_BasicScreen(t *testing.T) {
	name := "home"

	node := NewBuilder().
		Name(name).
		Update(func(*state.UIState, ScreenEvent) Result {
			return Result{}
		}).
		View(func(state.UIState) viewmodel.ViewModel {
			return viewmodel.ViewModel{}
		}).
		ToNode()

	assert.Equal(t, name, node.Screen.Name)
	assert.Len(t, 0, node.Stack)
	assert.Nil(t, node.Screen.Definition)
	assert.NotNil(t, node.Screen.Update)
	assert.NotNil(t, node.Screen.View)
}

func TestBuilder_WithoutDefinition(t *testing.T) {
	node := NewBuilder().
		Name("home").
		WithoutDefinition().
		ToNode()

	assert.NotNil(t, node.Screen.Definition)
	assert.Len(t, 0, node.Screen.Definition().RequireKeys)
}

func TestBuilder_NameToStack(t *testing.T) {
	name := "home"

	node := NewBuilder().
		Name(name).
		NameToStack().
		ToNode()

	assert.Contains(t, node.Stack, name)
}

func TestBuilder_IncompleteScreen(t *testing.T) {
	node := NewBuilder().Name("home").ToNode()

	assert.Nil(t, node.Screen.Definition)
	assert.Nil(t, node.Screen.Update)
	assert.Nil(t, node.Screen.View)
}
