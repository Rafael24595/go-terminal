package primitive

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestModalMenu_ToScreen(t *testing.T) {
	menu := NewModalMenu().
		SetName("base")

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestModalMenu_Stack(t *testing.T) {
	stack := NewModalMenu().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(default_modal_menu_name))
}
