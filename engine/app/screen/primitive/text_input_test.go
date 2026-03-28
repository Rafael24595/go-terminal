package primitive

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestTextInput_ToScreen(t *testing.T) {
	menu := NewTextInput().
		SetName("base")

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestTextInpu_Stack(t *testing.T) {
	stack := NewTextInput().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(default_text_input_name))
}
