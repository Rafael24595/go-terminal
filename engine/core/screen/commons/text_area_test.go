package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTextArea_ToScreen(t *testing.T) {
	menu := NewTextArea().
		SetName("base").
		AddTitle(text.LineFromString("Welcome"))

	screen := menu.ToScreen()

	Helper_ToScreen(t, screen)

	assert.Equal(t, screen.Name(), "base")
}
