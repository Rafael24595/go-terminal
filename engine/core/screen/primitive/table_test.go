package primitive

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/test/support/assert"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/core/screen"
)

func TestTable_ToScreen(t *testing.T) {
	menu := NewTable[int]().
		SetName("base").
		AddTitle(text.LineFromString("Welcome"))

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestTable_Stack(t *testing.T) {
	stack := NewTable[int]().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(default_table_name))
}
