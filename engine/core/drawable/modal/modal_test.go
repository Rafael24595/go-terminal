package modal

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/text"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestModal_ToDrawable(t *testing.T) {
	dw := ModalDrawableFromData(
		[]text.Line{},
		[]text.Fragment{},
		0,
	)

	drawable_test.Helper_ToDrawable(t, dw)
}

func TestModalDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	bd := NewModalDrawable()

	assert.Panic(t, func() {
		bd.draw()
	})
}
