package loop

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestLoop_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := LoopDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestLoop_Child_WipeCalled(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("golang"),
		},
	}

	dw := LoopDrawableFromDrawable(mock.ToDrawable())
	dw.Init()

	assert.False(t, mock.WipeCalled)

	dw.Draw(winsize.Winsize{})

	assert.True(t, mock.WipeCalled)
}

func TestLoopDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewLoopDrawable(mock.ToDrawable())

	assert.Panic(t, func() {
		bd.draw(winsize.Winsize{})
	})
}
