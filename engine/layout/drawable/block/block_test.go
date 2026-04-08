package block

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestBlock_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := BlockDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestBlockDrawable_Init_ShouldPropagateToChild(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBlockDrawable(mock.ToDrawable())

	bd.init()

	assert.True(t, bd.loaded)
	assert.True(t, mock.InitCalled)
}

func TestBlockDrawable_Draw_ShouldReturnEmptyIfRowsIsZero(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBlockDrawable(mock.ToDrawable())

	bd.init()

	lines, hasNext := bd.draw(terminal.Winsize{Rows: 0, Cols: 10})

	assert.Len(t, 0, lines)
	assert.True(t, hasNext)
}

func TestBlockDrawable_Draw_ShouldStopWhenChildHasNoNext(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{text.NewLine("golang")},
	}

	bd := NewBlockDrawable(mock.ToDrawable())
	bd.init()

	lines, hasNext := bd.draw(terminal.Winsize{Rows: 5, Cols: 10})

	assert.Len(t, 1, lines)
	assert.False(t, hasNext)
}

func TestBlockDrawable_Draw_ShouldAccumulateLines(t *testing.T) {
	count := 0

	mc := &drawable_test.MockDrawable{}
	dw := mc.ToDrawable()

	dw.Draw = func(size terminal.Winsize) ([]text.Line, bool) {
		count++
		return []text.Line{text.NewLine("golang")}, true
	}

	rows := uint16(3)
	bd := NewBlockDrawable(dw)
	bd.init()

	lines, hasNext := bd.draw(terminal.Winsize{Rows: rows, Cols: 10})

	assert.Len(t, int(rows), lines)
	assert.Equal(t, count, int(rows))
	assert.True(t, hasNext)
}
