package table

import (
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/commons/structure"
	"github.com/Rafael24595/go-terminal/engine/core"
	drawable_line "github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Use as a argument.
const min_width = 4

type section struct {
	header core.Drawable
	rows   core.Drawable
	footer core.Drawable
}

type col struct {
	name string
	size int
}

func makeSections(t table.Table, size terminal.Winsize) []section {
	sections := make([]section, 0)

	cols := int(size.Cols)
	separator := t.GetSeparator()
	headers := t.GetHeaders()

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
		headers := headersFromSize(table, headers)

		capacity := renderedRowSize(table, separator)
		specCover := style.SpecRepeatRight(uint(capacity))

		top := core.LineFromFragments(
			core.NewFragment(separator.Top).AddSpec(specCover),
		)

		bottom := core.LineFromFragments(
			core.NewFragment(separator.Bottom).AddSpec(specCover),
		)

		rows := makeTable(table, headers, t.GetColumns(), separator)
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

func headersFromSize(size map[string]int, headers []string) []string {
	filtered := make([]string, 0)
	for _, header := range headers {
		if _, ok := size[header]; ok {
			filtered = append(filtered, header)
		}
	}
	return filtered
}

func makeHeaders(size map[string]int, headers []string, separator table.SeparatorMeta) core.Line {
	headersLen := len(headers)

	capacity := 2*headersLen + 1
	fragments := make([]core.Fragment, 0, capacity)

	fragments = append(fragments, core.NewFragment(separator.Left))

	for i, h := range headers {
		width := uint(size[h])
		spec := style.MergeSpec(
			style.SpecPaddingCenter(width),
			style.SpecTrimRight(width),
		)
		fragments = append(fragments, core.NewFragment(h).AddSpec(spec))

		if i < headersLen-1 {
			fragments = append(fragments, core.NewFragment(separator.Center))
		}
	}

	fragments = append(fragments, core.NewFragment(separator.Right))

	return core.LineFromFragments(fragments...)
}

func makeTable(size map[string]int, headers []string, cols map[string][]string, separator table.SeparatorMeta) []core.Line {
	headersLen := len(headers)

	capacity := 2*headersLen + 1
	colSize := table.Cols(headers, cols)

	lines := make([]core.Line, colSize)

	for y := range colSize {
		fragments := make([]core.Fragment, 0, capacity)
		fragments = append(fragments, core.NewFragment(separator.Left))

		for i, h := range headers {
			width := uint(size[h])
			col := cols[h]

			if y >= 0 && y < len(col) {
				spec := style.MergeSpec(
					style.SpecPaddingRight(width),
					style.SpecTrimRight(width),
				)
				fragments = append(fragments, core.NewFragment(col[y]).AddSpec(spec))
			} else {
				spec := style.SpecRepeatRight(width)
				fragments = append(fragments, core.NewFragment("").AddSpec(spec))
			}

			if i < headersLen-1 {
				fragments = append(fragments, core.NewFragment(separator.Center))
			}
		}

		fragments = append(fragments, core.NewFragment(separator.Right))

		lines[y] = core.LineFromFragments(fragments...)
	}

	return lines
}

func renderedRowSize(size map[string]int, separator table.SeparatorMeta) int {
	sepCenterLen := utf8.RuneCountInString(separator.Center)
	sepLeftLen := utf8.RuneCountInString(separator.Left)
	sepRightLen := utf8.RuneCountInString(separator.Right)

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

	h := structure.NewMaxHeapBy(func(c col) int {
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

func splitTable(size map[string]int, headers []string, splitTable table.SeparatorMeta, cols int) []map[string]int {
	tables := make([]map[string]int, 0)

	leftLen := utf8.RuneCountInString(splitTable.Left)
	centerLen := utf8.RuneCountInString(splitTable.Center)
	rightLen := utf8.RuneCountInString(splitTable.Right)
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
