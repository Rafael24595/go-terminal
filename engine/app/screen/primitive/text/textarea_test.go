package text

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTextArea_ToScreen(t *testing.T) {
	menu := NewArea().
		SetName("base").
		AddTitle(*text.NewLine("Welcome"))

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestTextArea_Stack(t *testing.T) {
	stack := NewArea().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(area_name))
}
