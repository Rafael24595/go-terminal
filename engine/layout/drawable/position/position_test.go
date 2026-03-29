package position

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestPosition_ToDrawable(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := PositionDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestPositionDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewPositionDrawable(mock.ToDrawable())

	assert.Panic(t, func() {
		bd.draw()
	})
}
