package table

import (
	"fmt"
	"slices"
)

type SeparatorMeta struct {
	Top    string
	Bottom string
	Center string
	Left   string
	Right  string
}

var default_separator = SeparatorMeta{
	Top:    "-",
	Bottom: "-",
	Center: " | ",
	Left:   "| ",
	Right:  " |",
}

type Table struct {
	cols      map[string][]string
	headers   []string
	separator SeparatorMeta
}

func NewTable() *Table {
	return &Table{
		headers:   make([]string, 0),
		cols:      make(map[string][]string),
		separator: default_separator,
	}
}

func (t *Table) GetSeparator() SeparatorMeta {
	return t.separator
}

func (t *Table) SetSeparator(separator SeparatorMeta) *Table {
	t.separator = separator
	return t
}

func (t *Table) GetHeaders(headers ...string) []string {
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

func (t *Table) GetColumns(headers ...string) map[string][]string {
	return t.cols
}

func (t *Table) Field(header string, row int, data any) *Table {
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
	colSize := 0
	for _, h := range t.headers {
		if len(t.cols[h]) > colSize {
			colSize = len(t.cols[h])
		}
	}
	return colSize
}
