package table

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTable_ToScreen(t *testing.T) {
	menu := New[int]().
		SetName("base").
		AddTitle(*text.NewLine("Welcome"))

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestTable_Stack(t *testing.T) {
	stack := New[int]().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(name))
}
