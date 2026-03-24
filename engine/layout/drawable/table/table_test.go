package table

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/table"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestTable_ToDrawable(t *testing.T) {
	dw := TableDrawableFromTable(*table.NewTable(), *input.NewMatrixCursor(0, 0, false), style.Right)
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestTableDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	td := NewTableDrawable(*table.NewTable(), *input.NewMatrixCursor(0, 0, false), style.Right)

	assert.Panic(t, func() {
		td.draw()
	})
}
