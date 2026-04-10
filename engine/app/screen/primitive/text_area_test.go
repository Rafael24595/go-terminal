package primitive

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-terminal/engine/render/text"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestTextArea_ToScreen(t *testing.T) {
	menu := NewTextArea().
		SetName("base").
		AddTitle(*text.NewLine("Welcome"))

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestTextArea_Stack(t *testing.T) {
	stack := NewTextArea().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(default_text_area_name))
}
