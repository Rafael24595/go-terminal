package table

import (
	"fmt"
	"slices"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
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

func (t *Table) SetCell(header string, row uint16, data any) *Table {
	col, ok := t.cols[header]
	if !ok {
		return t
	}

	colLen := uint16(len(col))
	if row >= colLen {
		for i := colLen; i <= row; i++ {
			col = append(col, "")
		}
	}

	col[row] = fmt.Sprintf("%v", data)
	t.cols[header] = col

	return t
}

func (t *Table) Size() map[string]winsize.Cols {
	size := make(map[string]winsize.Cols)
	for _, h := range t.headers {
		if _, ok := size[h]; !ok {
			size[h] = runes.Measure(h)
		}

		for _, c := range t.cols[h] {
			size[h] = max(size[h], runes.Measure(c))
		}
	}

	return size
}

func (t *Table) Cols() uint16 {
	return uint16(len(t.headers))
}

func (t *Table) Rows() uint16 {
	return Rows(t.headers, t.cols)
}
