package static

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestStatic_ToDrawable(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := StaticDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestStaticDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	
	bd := NewStaticDrawable(mock.ToDrawable())
	bd.draw()

	bd = NewStaticDrawable(line.EagerLoopDrawableFromLines())
	assert.Panic(t, func() {
		bd.draw()
	})
}
