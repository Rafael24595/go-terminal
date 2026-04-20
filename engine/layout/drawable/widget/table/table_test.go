package table

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/table"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestTable_DrawableBasicSuite(t *testing.T) {
	dw := TableDrawableFromTable(
		*table.NewTable(),
		*input.NewMatrixCursor(0, 0, false),
	)
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestTable_LazyInit(t *testing.T) {
	dw := NewTableDrawable(
		*table.NewTable().
			SetHeaders("lang").
			SetCell("lang", 0, "golang"),
		*input.NewMatrixCursor(0, 0, false),
	)

	assert.Len(t, 0, dw.sections)

	dw.init()
	dw.draw(winsize.Winsize{
		Rows: 3,
		Cols: 11,
	})

	assert.Len(t, 1, dw.sections)
}
