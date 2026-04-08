package loop

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestLoop_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := LoopDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestLoop_Child_WipeCalled(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			text.NewLine("golang"),
		},
	}

	dw := LoopDrawableFromDrawable(mock.ToDrawable())
	dw.Init()

	assert.False(t, mock.WipeCalled)

	dw.Draw(terminal.Winsize{})

	assert.True(t, mock.WipeCalled)
}

func TestLoopDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewLoopDrawable(mock.ToDrawable())

	assert.Panic(t, func() {
		bd.draw(terminal.Winsize{})
	})
}
