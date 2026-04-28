package text

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTextInput_ToScreen(t *testing.T) {
	menu := NewInput().
		SetName("base")

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestTextInput_Stack(t *testing.T) {
	stack := NewInput().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(input_name))
}
