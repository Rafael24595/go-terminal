package table

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/table"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTable_ToDrawable(t *testing.T) {
	dw := TableDrawableFromTable(*table.NewTable(), *NewCursor(0, 0, false), Right)
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestTableDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	td := NewTableDrawable(*table.NewTable(), *NewCursor(0, 0, false), Right)

	assert.Panic(t, func() {
		td.draw()
	})
}
