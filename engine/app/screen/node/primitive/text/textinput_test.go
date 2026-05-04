package text

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTextInput_ToNode(t *testing.T) {
	node := NewInput().SetName("base").ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Screen.Name, "base")
}

func TestTextInput_Stack(t *testing.T) {
	stack := NewInput().ToNode().Stack

	assert.True(t, stack.Has(NameInput))
}
