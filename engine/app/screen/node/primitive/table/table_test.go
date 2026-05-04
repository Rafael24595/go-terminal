package table

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTable_ToNode(t *testing.T) {
	node := New[int]().
		SetName("base").
		AddTitle(*text.NewLine("Welcome")).
		ToNode()

	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Screen.Name, "base")
}

func TestTable_Stack(t *testing.T) {
	stack := New[int]().ToNode().Stack

	assert.True(t, stack.Has(Name))
}
