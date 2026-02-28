package box

import (
	"testing"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTable_ToDrawable(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := BoxDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestTableDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBoxDrawable(mock.ToDrawable())

	assert.Panic(t, func() {
		bd.draw()
	})
}
