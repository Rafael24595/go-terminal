package textarea

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/input"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTextArea_ToDrawable(t *testing.T) {
	dw := TextAreaDrawableFromData([]rune{}, *input.NewTextCursor(false))
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestTextAreaDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	td := NewTextAreaDrawable([]rune{}, *input.NewTextCursor(false))

	assert.Panic(t, func() {
		td.draw()
	})
}
