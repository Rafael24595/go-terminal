package wrapper

import (
	"testing"

	"github.com/Rafael24595/go-terminal/test/support/assert"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/core/screen"
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
