package primitive

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/test/support/mock"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestCheckMenu_ToScreen(t *testing.T) {
	menu := NewCheckMenu().
		SetName("base").
		AddTitle(*text.NewLine("Welcome"))

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestCheckMenu_Stack(t *testing.T) {
	stack := NewCheckMenu().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(default_check_menu_name))
}

func TestCheckMenu_SwitchState_WithLimit(t *testing.T) {
	clock := &mock.TestClock{Time: 1000}

	menu := NewCheckMenu().
		SetLimit(2).
		AddOptions(
			input.NewCheckOption("1", *text.NewFragment("option 1")),
			input.NewCheckOption("2", *text.NewFragment("option 2")),
			input.NewCheckOption("3", *text.NewFragment("option 3")),
		)

	menu.clock = clock.Now

	assert.False(t, menu.options[0].Status)
	assert.False(t, menu.options[1].Status)
	assert.False(t, menu.options[2].Status)

	clock.Advance(1000)

	menu.cursor = 0
	menu.switchState()

	assert.True(t, menu.options[0].Status)
	assert.False(t, menu.options[1].Status)
	assert.False(t, menu.options[2].Status)

	clock.Advance(1000)

	menu.cursor = 1
	menu.switchState()

	assert.True(t, menu.options[0].Status)
	assert.True(t, menu.options[1].Status)
	assert.False(t, menu.options[2].Status)

	clock.Advance(1000)

	menu.cursor = 2
	menu.switchState()

	assert.False(t, menu.options[0].Status)
	assert.True(t, menu.options[1].Status)
	assert.True(t, menu.options[2].Status)
}
