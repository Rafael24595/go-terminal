package primitive

import (
	"testing"

	"github.com/Rafael24595/go-terminal/test/support/assert"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/core/screen"
)

func TestModalMenu_ToScreen(t *testing.T) {
	menu := NewModalMenu().
		SetName("base")

	screen := menu.ToScreen()

	screen_test.Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}
