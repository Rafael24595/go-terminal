package form

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/config/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestForm_ToNode(t *testing.T) {
	node := New().ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, Name, node.Name)
}

func TestForm_Propagate(t *testing.T) {
	name := "base"
	mock := screen_test.MockScreen{
		Name: name,
	}

	node := New().
		AddNode(
			true,
			mock.ToNode(),
			chunk.Chunk[winsize.Rows]{},
		).
		ToNode()

	screen_test.Helper_Propagate(t, name, 0, node)
}
