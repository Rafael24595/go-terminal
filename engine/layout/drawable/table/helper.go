package table

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure/heap"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	drawable_line "github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/table"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Use as a argument.
const min_width = int(3 + marker.DefaultElipsisSize)

type section struct {
	header drawable.Drawable
	rows   drawable.Drawable
	footer drawable.Drawable
}

type col struct {
	name string
	size int
}

func makeSections(t table.Table, cursor input.MatrixCursor, size terminal.Winsize) []section {
	sections := make([]section, 0)

	cols := int(size.Cols)
	separator := t.GetSeparator()
	headers := t.GetHeaders()
	columns := t.GetColumns()

	baseSize := t.Size()
	rendSize := renderedRowSize(baseSize, separator)
	realSize, status := adjustSize(baseSize, headers, cols, rendSize)

	var tables []map[string]int
	if !status {
		tables = splitTable(realSize, headers, separator, cols)
	} else {
		tables = []map[string]int{realSize}
	}

	for _, table := range tables {
		headers, fixCursor := headersFromSize(table, headers, cursor)

		capacity := renderedRowSize(table, separator)
		specCover := style.SpecRepeatRight(uint(capacity))

		top := text.LineFromFragments(
			text.NewFragment(separator.Top).AddSpec(specCover),
		)

		bottom := text.LineFromFragments(
			text.NewFragment(separator.Bottom).AddSpec(specCover),
		)

		rows := makeTable(table, headers, columns, separator, fixCursor)

		if len(rows) == 0 {
			continue
		}

		headerRow := makeHeaders(table, headers, separator)

		sections = append(sections, section{
			header: drawable_line.EagerLoopDrawableFromLines(
				top, headerRow, top,
			),
			rows:   drawable_line.LazyDrawableFromLines(rows...),
			footer: drawable_line.EagerLoopDrawableFromLines(bottom),
		})
	}

	return sections
}

func headersFromSize(size map[string]int, headers []string, cursor input.MatrixCursor) ([]string, *input.MatrixCursor) {
	filtered := make([]string, 0)

	var fixCursor *input.MatrixCursor

	fixX := 0
	for x, header := range headers {
		if _, ok := size[header]; !ok {
			continue
		}

		filtered = append(filtered, header)
		if x == int(cursor.Col) {
			fixCursor = input.NewMatrixCursor(
				cursor.Row,
				uint32(fixX),
				cursor.Show,
			)
		}
		fixX += 1
	}

	return filtered, fixCursor
}

func makeHeaders(size map[string]int, headers []string, separator marker.TableSeparatorMeta) text.Line {
	headersLen := len(headers)

	capacity := 2*headersLen + 1
	fragments := make([]text.Fragment, 0, capacity)

	fragments = append(fragments, text.NewFragment(separator.Left))

	for i, h := range headers {
		width := uint(size[h])
		spec := style.MergeSpec(
			style.SpecPaddingCenter(width),
			style.SpecTrimTextRight(width, marker.DefaultElipsisText),
		)
		fragments = append(fragments, text.NewFragment(h).AddSpec(spec))

		if i < headersLen-1 {
			fragments = append(fragments, text.NewFragment(separator.Center))
		}
	}

	fragments = append(fragments, text.NewFragment(separator.Right))

	return text.LineFromFragments(fragments...)
}

func makeTable(
	size map[string]int,
	headers []string,
	cols map[string][]string,
	separator marker.TableSeparatorMeta,
	cursor *input.MatrixCursor,
) []text.Line {
	headersLen := len(headers)

	capacity := 2*headersLen + 1
	maxRow := table.Rows(headers, cols)

	lines := make([]text.Line, maxRow)

	for y := range maxRow {
		fragments := make([]text.Fragment, 0, capacity)
		fragments = append(fragments, text.NewFragment(separator.Left))

		for x, h := range headers {
			frag := makeCell(size, cols, cursor, h, y, x)
			fragments = append(fragments, frag)

			if x < headersLen-1 {
				fragments = append(fragments, text.NewFragment(separator.Center))
			}
		}

		fragments = append(fragments, text.NewFragment(separator.Right))

		lines[y] = text.LineFromFragments(fragments...)
	}

	return lines
}

func makeCell(
	size map[string]int,
	cols map[string][]string,
	cursor *input.MatrixCursor,
	header string,
	y int,
	x int,
) text.Fragment {
	width := uint(size[header])
	col := cols[header]

	atom := style.AtmNone

	cursorShow := cursor != nil && cursor.Show
	if cursorShow && y == int(cursor.Row) && x == int(cursor.Col) {
		atom = style.AtmSelect | style.AtmFocus
	}

	if y >= 0 && y < len(col) {
		spec := style.MergeSpec(
			style.SpecPaddingRight(width),
			style.SpecTrimTextRight(width, marker.DefaultElipsisText),
		)

		return text.NewFragment(col[y]).
			AddSpec(spec).
			AddAtom(atom)
	}

	spec := style.SpecRepeatRight(width)

	return text.NewFragment("").
		AddSpec(spec).
		AddAtom(atom)
}

func renderedRowSize(size map[string]int, separator marker.TableSeparatorMeta) int {
	sepCenterLen := runes.Measure(separator.Center)
	sepLeftLen := runes.Measure(separator.Left)
	sepRightLen := runes.Measure(separator.Right)

	joinSize := (len(size) - 1) * sepCenterLen
	borderSize := sepLeftLen + sepRightLen

	total := 0
	for _, v := range size {
		total += v
	}

	return total + joinSize + borderSize
}

func adjustSize(size map[string]int, headers []string, cols int, rowSize int) (map[string]int, bool) {
	if rowSize <= cols {
		return size, true
	}

	excess := rowSize - cols

	h := heap.NewMaxHeapBy(func(c col) int {
		return c.size
	})

	for _, v := range headers {
		h.Push(col{v, size[v]})
	}

	for excess > 0 {
		c, ok := h.Peek()
		if !ok || c.size <= min_width {
			break
		}

		c, _ = h.Pop()
		c.size--
		excess--
		h.Push(c)
	}

	newSize := make(map[string]int)
	for h.Len() > 0 {
		c, _ := h.Pop()
		newSize[c.name] = c.size
	}

	return newSize, excess == 0
}

func splitTable(size map[string]int, headers []string, splitTable marker.TableSeparatorMeta, cols int) []map[string]int {
	tables := make([]map[string]int, 0)

	leftLen := runes.Measure(splitTable.Left)
	centerLen := runes.Measure(splitTable.Center)
	rightLen := runes.Measure(splitTable.Right)
	headersLen := len(headers)

	table := make(map[string]int)
	count := leftLen

	for i, k := range headers {
		v := min(size[k], cols)

		needed := count + v
		if len(table) > 0 {
			needed += centerLen
		}

		needed += rightLen

		if needed > cols && len(table) > 0 {
			tables = append(tables, table)

			table = make(map[string]int)
			count = 0
		}

		table[k] = v
		count += v

		if i < headersLen-1 {
			count += centerLen
		}
	}

	if len(table) != 0 {
		tables = append(tables, table)
	}

	return tables
}
