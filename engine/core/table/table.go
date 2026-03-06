package table

import (
	"fmt"
	"slices"

	"github.com/Rafael24595/go-terminal/engine/core/marker"
)

type Table struct {
	cols      map[string][]string
	headers   []string
	separator marker.TableSeparatorMeta
}

func NewTable() *Table {
	return &Table{
		headers:   make([]string, 0),
		cols:      make(map[string][]string),
		separator: marker.DefaultTableSeparator,
	}
}

func (t *Table) GetSeparator() marker.TableSeparatorMeta {
	return t.separator
}

func (t *Table) SetSeparator(separator marker.TableSeparatorMeta) *Table {
	t.separator = separator
	return t
}

func (t *Table) GetHeaders() []string {
	return t.headers
}

func (t *Table) SetHeaders(headers ...string) *Table {
	for _, v := range headers {
		if slices.Contains(t.headers, v) {
			continue
		}

		t.headers = append(t.headers, v)
		t.cols[v] = make([]string, 0)
	}

	return t
}

func (t *Table) GetColumns() map[string][]string {
	return t.cols
}

func (t *Table) FindCellByCoords(row, col int) (string, bool) {
	if col >= len(t.headers) {
		return "", false
	}

	header := t.headers[col]

	cols, ok := t.cols[header]
	if !ok || row > len(cols) {
		return "", false
	}

	return cols[row], true
}

func (t *Table) SetCell(header string, row int, data any) *Table {
	col, ok := t.cols[header]
	if !ok {
		return t
	}

	if row >= len(col) {
		for i := len(col); i <= row; i++ {
			col = append(col, "")
		}
	}

	col[row] = fmt.Sprintf("%v", data)
	t.cols[header] = col

	return t
}

func (t *Table) Size() map[string]int {
	size := make(map[string]int)
	for _, h := range t.headers {
		if _, ok := size[h]; !ok {
			size[h] = len(h)
		}

		for _, c := range t.cols[h] {
			size[h] = max(size[h], len(c))
		}
	}

	return size
}

func (t *Table) Cols() int {
	return len(t.headers)
}

func (t *Table) Rows() int {
	return Rows(t.headers, t.cols)
}
