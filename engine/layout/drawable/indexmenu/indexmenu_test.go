package indexmenu

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestIndexMenu_ToDrawable(t *testing.T) {
	dw := TextIndexMenuFromData([]text.Fragment{})
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestIndexMenuDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	td := NewIndexMenuDrawable([]text.Fragment{})

	assert.Panic(t, func() {
		td.draw()
	})
}
