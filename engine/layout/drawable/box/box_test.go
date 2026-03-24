package box

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestBox_ToDrawable(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := BoxDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestBoxDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBoxDrawable(mock.ToDrawable())

	assert.Panic(t, func() {
		bd.draw()
	})
}
