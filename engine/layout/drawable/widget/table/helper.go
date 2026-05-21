package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/heap"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/isolated"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_line "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
)

// TODO: Use as a argument.
const min_width = 3 + marker.DefaultElipsisSize

type section struct {
	header drawable.Unit
	rows   drawable.Unit
	footer drawable.Unit
}

type col struct {
	name string
	size winsize.Cols
}

func makeSections(t table.Table, cursor input.MatrixCursor, size winsize.Winsize) []section {
	sections := make([]section, 0)

	cols := size.Cols
	separator := t.GetSeparator()
	headers := t.GetHeaders()
	columns := t.GetColumns()

	baseSize := t.Size()
	rendSize := renderedRowSize(baseSize, separator)
	realSize, status := adjustSize(baseSize, headers, cols, rendSize)

	var tables []map[string]winsize.Cols
	if !status {
		tables = splitTable(realSize, headers, separator, cols)
	} else {
		tables = []map[string]winsize.Cols{realSize}
	}

	for _, table := range tables {
		headers, fixCursor := headersFromSize(table, headers, cursor)

		capacity := renderedRowSize(table, separator)
		specCover := style.SpecRepeatRight(capacity)

		top := text.LineFromFragments(
			*text.NewFragment(separator.Top).AddSpec(specCover),
		)

		bottom := text.LineFromFragments(
			*text.NewFragment(separator.Bottom).AddSpec(specCover),
		)

		rows := makeTable(table, headers, columns, separator, fixCursor)

		if len(rows) == 0 {
			continue
		}

		headerRow := makeHeaders(table, headers, separator)

		sections = append(sections, section{
			header: isolated.UnitFromLines(
				*top, *headerRow, *top,
			),
			rows: drawable_line.UnitFromLines(rows...),
			footer: isolated.UnitFromLines(
				*bottom,
			),
		})
	}

	return sections
}

func headersFromSize(
	size map[string]winsize.Cols,
	headers []string,
	cursor input.MatrixCursor,
) ([]string, *input.MatrixCursor) {
	filtered := make([]string, 0)

	var fixCursor *input.MatrixCursor

	fixX := uint16(0)
	for x, header := range headers {
		if _, ok := size[header]; !ok {
			continue
		}

		filtered = append(filtered, header)
		if x == int(cursor.Col) {
			fixCursor = input.NewMatrixCursor(
				cursor.Row,
				fixX,
				cursor.Show,
			)
		}
		fixX += 1
	}

	return filtered, fixCursor
}

func makeHeaders(
	size map[string]winsize.Cols,
	headers []string,
	separator marker.TableSeparatorMeta,
) *text.Line {
	headersLen := len(headers)

	capacity := 2*headersLen + 1
	fragments := make([]text.Fragment, 0, capacity)

	fragments = append(fragments, *text.NewFragment(separator.Left))

	for i, h := range headers {
		width := size[h]
		spec := style.MergeSpec(
			style.SpecPaddingCenter(width),
			style.SpecTrimTextRight(width, marker.DefaultElipsisText),
		)
		fragments = append(fragments, *text.NewFragment(h).AddSpec(spec))

		if i < headersLen-1 {
			fragments = append(fragments, *text.NewFragment(separator.Center))
		}
	}

	fragments = append(fragments, *text.NewFragment(separator.Right))

	return text.LineFromFragments(fragments...)
}

func makeTable(
	size map[string]winsize.Cols,
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

		lSep := *text.NewFragment(separator.Left).
			AddAtom(style.AtmWrap)

		fragments = append(fragments, lSep)

		for x, h := range headers {
			frag := makeCell(size, cols, cursor, h, y, uint16(x))
			fragments = append(fragments, *frag)

			if x < headersLen-1 {
				cSep := *text.NewFragment(separator.Center).
					AddAtom(style.AtmWrap)

				fragments = append(fragments, cSep)
			}
		}

		rSep := *text.NewFragment(separator.Right).
			AddAtom(style.AtmWrap)

		fragments = append(fragments, rSep)

		lines[y] = *text.LineFromFragments(fragments...)
	}

	return lines
}

func makeCell(
	size map[string]winsize.Cols,
	cols map[string][]string,
	cursor *input.MatrixCursor,
	header string,
	y uint16,
	x uint16,
) *text.Fragment {
	width := size[header]
	col := cols[header]

	atom := style.AtmWrap

	cursorShow := cursor != nil && cursor.Show
	if cursorShow && y == cursor.Row && x == cursor.Col {
		atom = style.MergeAtom(atom, style.AtmSelect, style.AtmFocus)
	}

	if y < uint16(len(col)) {
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

func renderedRowSize(size map[string]winsize.Cols, separator marker.TableSeparatorMeta) winsize.Cols {
	sepCenterLen := runes.Measure(separator.Center)
	sepLeftLen := runes.Measure(separator.Left)
	sepRightLen := runes.Measure(separator.Right)

	mapLen := max(0, len(size)-1)
	joinSize := winsize.Cols(mapLen) * sepCenterLen
	borderSize := sepLeftLen + sepRightLen

	total := winsize.Cols(0)
	for _, v := range size {
		total += v
	}

	return total + joinSize + borderSize
}

func adjustSize(
	size map[string]winsize.Cols,
	headers []string,
	cols winsize.Cols,
	rowSize winsize.Cols,
) (map[string]winsize.Cols, bool) {
	if rowSize <= cols {
		return size, true
	}

	excess := rowSize.Sub(cols)

	h := heap.NewMaxHeapBy(func(c col) winsize.Cols {
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

		c.size = c.size.Sub(1)
		excess = excess.Sub(1)

		h.Push(c)
	}

	newSize := make(map[string]winsize.Cols)
	for h.Len() > 0 {
		c, _ := h.Pop()
		newSize[c.name] = c.size
	}

	return newSize, excess == 0
}

func splitTable(
	size map[string]winsize.Cols,
	headers []string,
	splitTable marker.TableSeparatorMeta,
	cols winsize.Cols,
) []map[string]winsize.Cols {
	tables := make([]map[string]winsize.Cols, 0)

	leftLen := runes.Measure(splitTable.Left)
	centerLen := runes.Measure(splitTable.Center)
	rightLen := runes.Measure(splitTable.Right)
	headersLen := len(headers)

	table := make(map[string]winsize.Cols)
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

			table = make(map[string]winsize.Cols)
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
