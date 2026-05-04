package modalmenu

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestModalMenu_ToNode(t *testing.T) {
	node := New().SetName("base").ToNode()

	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Screen.Name, "base")
}

func TestModalMenu_Stack(t *testing.T) {
	stack := New().ToNode().Stack

	assert.True(t, stack.Has(Name))
}
