package indexmenu

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTextArea_ToDrawable(t *testing.T) {
	dw := TextIndexMenuFromData(marker.HyphenIndex, []text.Fragment{})
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestTextAreaDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	td := NewIndexMenuDrawable(marker.HyphenIndex, []text.Fragment{})

	assert.Panic(t, func() {
		td.draw()
	})
}
