package box

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/terminal"
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

func TestBoxDrawable_Init_ShouldPropagateToChild(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBoxDrawable(mock.ToDrawable())

	size := terminal.Winsize{Rows: 10, Cols: 20}
	bd.init(size)

	assert.True(t, bd.initialized)
	assert.Equal(t, bd.size, size)
	assert.True(t, mock.InitCalled)
}
