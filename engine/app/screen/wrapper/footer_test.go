package wrapper

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestFooter_ToScreen(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	menu := NewFooter(mock.ToScreen())

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}

func TestFooter_Stack(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	stack := NewFooter(mock.ToScreen()).
		ToScreen().
		Stack()

	assert.True(t, stack.Has("base"))
}
