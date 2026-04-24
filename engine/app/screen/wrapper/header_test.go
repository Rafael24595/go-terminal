package wrapper

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestHeader_ToScreen(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	menu := NewHeader(mock.ToScreen())

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestHeader_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := NewHeader(mock.ToScreen()).
		ToScreen().
		Stack()

	assert.True(t, stack.Has("base"))
}
