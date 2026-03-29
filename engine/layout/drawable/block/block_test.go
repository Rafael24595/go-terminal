package block

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestBlock_ToDrawable(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := BlockDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestBlockDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBlockDrawable(mock.ToDrawable())

	assert.Panic(t, func() {
		bd.draw()
	})
}

func TestBlockDrawable_Init_ShouldPropagateToChild(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBlockDrawable(mock.ToDrawable())

	size := terminal.Winsize{Rows: 10, Cols: 20}
	bd.init(size)

	assert.True(t, bd.initialized)
	assert.Equal(t, bd.size, size)
	assert.True(t, mock.InitCalled)
}

func TestBlockDrawable_Draw_ShouldReturnEmptyIfRowsIsZero(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewBlockDrawable(mock.ToDrawable())

	bd.init(terminal.Winsize{Rows: 0, Cols: 10})

	lines, hasNext := bd.draw()

	assert.Equal(t, len(lines), 0)
	assert.True(t, hasNext)
}

func TestBlockDrawable_Draw_ShouldStopWhenChildHasNoNext(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{text.LineFromString("golang")},
	}

	bd := NewBlockDrawable(mock.ToDrawable())
	bd.init(terminal.Winsize{Rows: 5, Cols: 10})

	lines, hasNext := bd.draw()

	assert.Equal(t, len(lines), 1)
	assert.False(t, hasNext)
}

func TestBlockDrawable_Draw_ShouldAccumulateLines(t *testing.T) {
	count := 0

	mc := &drawable_test.MockDrawable{}
	dw := mc.ToDrawable()

	dw.Draw = func() ([]text.Line, bool) {
		count++
		return []text.Line{text.LineFromString("golang")}, true
	}

	rows := uint16(3)
	bd := NewBlockDrawable(dw)
	bd.init(terminal.Winsize{Rows: rows, Cols: 10})

	lines, hasNext := bd.draw()

	assert.Equal(t, len(lines), int(rows))
	assert.Equal(t, count, int(rows))
	assert.True(t, hasNext)
}
