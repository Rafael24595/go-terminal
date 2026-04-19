package box

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestBox_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := BoxDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestBoxDrawable_Init_ShouldPropagateToChild(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBoxDrawable(mock.ToDrawable())

	bd.init()

	assert.True(t, bd.loaded)
	assert.True(t, mock.InitCalled)
}
