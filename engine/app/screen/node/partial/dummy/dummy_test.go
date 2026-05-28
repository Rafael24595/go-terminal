package dummy

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestDummy_ToNode(t *testing.T) {
	node := ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, Name, node.Name)
}

func TestDummy_Defaults(t *testing.T) {
	node := ToNode()

	assert.Equal(t, Name, node.Name)
	assert.Len(t, 0, node.Tags)
	assert.Len(t, 0, node.Children())

	definition := node.Screen.Definition()
	assert.Equal(t, 0, definition.Descriptor.Size())
	assert.Equal(t, 0, definition.RequireKeys.Size())

	uiState := state.NewUIState()
	update := node.Screen.Update(
		uiState,
		screen.NewEvent(key.Key{}),
	)

	assert.Nil(t, update.Node)
	assert.False(t, update.Isolate)
	assert.DeepEqual(t, uiState.Pager, update.Pager)

	view := node.Screen.View(*uiState)
	assert.False(t, view.Behavior.NeedsPulse)
	assert.Equal(t, 0, view.Header.Size())
	assert.Equal(t, 0, view.Kernel.Size())
	assert.Equal(t, 0, view.Footer.Size())
}
