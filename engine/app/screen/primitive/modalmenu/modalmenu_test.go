package modalmenu

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestModalMenu_ToScreen(t *testing.T) {
	menu := New().
		SetName("base")

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestModalMenu_Stack(t *testing.T) {
	stack := New().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(name))
}
