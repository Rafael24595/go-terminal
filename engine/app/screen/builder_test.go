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
		Tick(func(*state.UIState, Event) Result {
			return Result{}
		}).
		View(func(state.UIState) viewmodel.ViewModel {
			return viewmodel.ViewModel{}
		}).
		ToNode()

	assert.Equal(t, name, node.Name)
	assert.Len(t, 0, node.Stack)
	assert.Nil(t, node.Screen.Keys)
	assert.NotNil(t, node.Screen.Tick)
	assert.NotNil(t, node.Screen.View)
}

func TestBuilder_WithoutKeys(t *testing.T) {
	node := NewBuilder().
		Name("home").
		WithoutKeys().
		ToNode()

	assert.NotNil(t, node.Screen.Keys)
	assert.Equal(t, 0, node.Screen.Keys().RequireKeys.Size())
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

	assert.Nil(t, node.Screen.Keys)
	assert.Nil(t, node.Screen.Tick)
	assert.Nil(t, node.Screen.View)
}
