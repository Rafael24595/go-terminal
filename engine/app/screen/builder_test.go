package screen

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

func TestBuilder_BasicScreen(t *testing.T) {
	name := "home"

	screen := NewBuilder().
		Name(name).
		Update(func(*state.UIState, ScreenEvent) Result {
			return Result{}
		}).
		View(func(state.UIState) viewmodel.ViewModel {
			return viewmodel.ViewModel{}
		}).
		ToScreen()

	assert.Equal(t, name, screen.Name)
	assert.Len(t, 0, screen.Stack)
	assert.Nil(t, screen.Definition)
	assert.NotNil(t, screen.Update)
	assert.NotNil(t, screen.View)
}

func TestBuilder_WithoutDefinition(t *testing.T) {
	screen := NewBuilder().
		Name("home").
		WithoutDefinition().
		ToScreen()

	assert.NotNil(t, screen.Definition)
	assert.Len(t, 0, screen.Definition().RequireKeys)
}

func TestBuilder_NameToStack(t *testing.T) {
	name := "home"

	screen := NewBuilder().
		Name(name).
		NameToStack().
		ToScreen()

	assert.Contains(t, screen.Stack, name)
}

func TestBuilder_IncompleteScreen(t *testing.T) {
	screen := NewBuilder().
		Name("home").
		ToScreen()

	assert.Nil(t, screen.Definition)
	assert.Nil(t, screen.Update)
	assert.Nil(t, screen.View)
}
