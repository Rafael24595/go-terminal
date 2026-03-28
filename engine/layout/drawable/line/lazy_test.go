package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestLazy_ToDrawable(t *testing.T) {
	dw := LazyDrawableFromLines()
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestLazyDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	bd := NewLazyDrawable()

	assert.Panic(t, func() {
		bd.draw()
	})
}
