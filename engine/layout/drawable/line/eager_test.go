package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestEager_ToDrawable(t *testing.T) {
	dw := EagerDrawableFromLines()
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestEagerDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	bd := NewEagerDrawable()

	assert.Panic(t, func() {
		bd.draw()
	})
}
