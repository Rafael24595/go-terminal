package help

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/model/help"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestHelp_ToDrawable(t *testing.T) {
	dw := HelpDrawableFromMeta(help.NewHelpMeta())
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestHelpDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	bd := NewHelpDrawable(help.NewHelpMeta())

	assert.Panic(t, func() {
		bd.draw()
	})
}
