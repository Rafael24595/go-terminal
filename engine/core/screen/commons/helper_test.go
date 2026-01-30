package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func Helper_ToScreen(t *testing.T, screen screen.Screen) {
	t.Helper()

	assert.NotNil(t, screen.Name, "Screen.Name")
	assert.NotNil(t, screen.View, "Screen.View should be set")
	assert.NotNil(t, screen.Update, "Screen.Update should be set")
}
