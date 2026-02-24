package table

import (
	"testing"

	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestNewTable_ShouldInitializeEmptyTable(t *testing.T) {
	tbl := NewTable()

	assert.Equal(t, 0, tbl.Cols())
	assert.Equal(t, 0, tbl.Rows())
	assert.Equal(t, default_separator, tbl.GetSeparator())
}

func TestSetHeaders_ShouldAddHeadersWithoutDuplicates(t *testing.T) {
	tbl := NewTable()

	tbl.SetHeaders("ID", "Lang")
	tbl.SetHeaders("Lang", "Age")

	headers := tbl.GetHeaders()

	assert.Len(t, 3, headers)

	assert.Equal(t, "ID", headers[0])
	assert.Equal(t, "Lang", headers[1])
	assert.Equal(t, "Age", headers[2])
}

func TestField_ShouldExpandRowsDynamically(t *testing.T) {
	tbl := NewTable()
	tbl.SetHeaders("Name")

	tbl.Field("Name", 2, "Golang")

	col := tbl.GetColumns()["Name"]

	assert.Len(t, 3, col)
	assert.Equal(t, "", col[0])
	assert.Equal(t, "", col[1])
	assert.Equal(t, "Golang", col[2])
}

func TestField_WithInvalidHeader_ShouldDoNothing(t *testing.T) {
	tbl := NewTable()
	tbl.SetHeaders("ID")

	tbl.Field("Invalid", 0, "X")

	assert.Len(t, 0, tbl.GetColumns()["ID"])
}

func TestSize_ShouldCalculateMaxWidth(t *testing.T) {
	tbl := NewTable()
	tbl.SetHeaders("Name")

	tbl.Field("Name", 0, "zig")
	tbl.Field("Name", 1, "golang")

	size := tbl.Size()

	assert.Len(t, size["Name"], []rune("golang"))
}

func TestSize_ShouldConsiderHeaderLength(t *testing.T) {
	tbl := NewTable()
	tbl.SetHeaders("VeryLongHeader")

	tbl.Field("VeryLongHeader", 0, "go")

	size := tbl.Size()

	assert.Len(t, size["VeryLongHeader"], []rune("VeryLongHeader"))
}

func TestCols_ShouldReturnHeaderCount(t *testing.T) {
	tbl := NewTable()
	tbl.SetHeaders("A", "B", "C")

	assert.Equal(t, 3, tbl.Cols())
}

func TestRows_ShouldReturnMaxRowCount(t *testing.T) {
	tbl := NewTable()
	tbl.SetHeaders("A", "B")

	tbl.Field("A", 0, "x")
	tbl.Field("B", 2, "y")

	assert.Equal(t, 3, tbl.Rows())
}

func TestSetSeparator_ShouldOverrideDefault(t *testing.T) {
	tbl := NewTable()

	sep := SeparatorMeta{
		Top:    "=",
		Bottom: "=",
		Center: "::",
		Left:   "[",
		Right:  "]",
	}

	ret := tbl.SetSeparator(sep)

	assert.Equal(t, sep, tbl.GetSeparator())
	assert.Equal(t, ret, tbl)
}
