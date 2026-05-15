package text

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTextArea_ToNode(t *testing.T) {
	node := NewArea().SetName("base").ToNode()

	screen_test.Helper_ToNode(t, node)
	assert.Equal(t, node.Screen.Name, "base")
}

func TestTextArea_Stack(t *testing.T) {
	stack := NewArea().ToNode().Stack

	assert.True(t, stack.Has(NameArea))
}
