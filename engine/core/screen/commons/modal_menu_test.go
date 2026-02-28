package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestModalMenu_ToScreen(t *testing.T) {
	menu := NewModalMenu().
		SetName("base")

	screen := menu.ToScreen()

	Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}
