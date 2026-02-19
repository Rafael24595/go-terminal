package table

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure"
	"github.com/Rafael24595/go-terminal/engine/core"
	drawable_line "github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Use as a argument.
const min_width = 4

type col struct {
	name string
	size int
}

func makeDrawables(t table.Table, size terminal.Winsize) []core.Drawable {
	drawables := make([]core.Drawable, 0)

	cols := int(size.Cols)
	separator := t.GetSeparator()

	baseSize := t.Size()
	rendSize := renderedRowSize(baseSize, separator)
	realSize, status := adjustSize(baseSize, cols, rendSize)

	var tables []map[string]int
	if !status {
		tables = splitTable(realSize, cols)
	} else {
		tables = []map[string]int{realSize}
	}

	for _, table := range tables {
		lines := make([]core.Line, 0)

		headers := headersFromSize(table, t.GetHeaders())

		capacity := renderedRowSize(table, separator)
		specCover := style.SpecRepeatRight(uint(capacity))

		top := core.LineFromFragments(
			core.NewFragment(separator.Top).AddSpec(specCover),
		)
		lines = append(lines, top)

		headerRow := makeHeaders(table, headers, separator)
		lines = append(lines, headerRow)
		lines = append(lines, top)

		rows := makeTable(table, headers, t.GetColumns(), separator)
		lines = append(lines, rows...)

		bottom := core.LineFromFragments(
			core.NewFragment(separator.Bottom).AddSpec(specCover),
		)
		lines = append(lines, bottom)

		if len(lines) > 0 {
			drawable := drawable_line.LinesEagerDrawableFromLines(lines...)
			drawables = append(drawables, drawable)
		}
	}

	return drawables
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
	capacity := 2*len(headers) + 1
	fragments := make([]core.Fragment, 0, capacity)

	fragments = append(fragments, core.NewFragment(separator.Left))

	for i, h := range headers {
		spec := style.SpecPaddingCenter(uint(size[h]))
		fragments = append(fragments, core.NewFragment(h).AddSpec(spec))

		if i < len(headers)-1 {
			fragments = append(fragments, core.NewFragment(separator.Center))
		}
	}

	fragments = append(fragments, core.NewFragment(separator.Right))

	return core.LineFromFragments(fragments...)
}

func makeTable(size map[string]int, headers []string, cols map[string][]string, separator table.SeparatorMeta) []core.Line {
	capacity := 2*len(headers) + 1
	colSize := table.Cols(headers, cols)

	lines := make([]core.Line, colSize)

	for y := range colSize {
		fragments := make([]core.Fragment, 0, capacity)
		fragments = append(fragments, core.NewFragment(separator.Left))

		for i, h := range headers {
			size := uint(size[h])
			col := cols[h]

			if y >= 0 && y < len(col) {
				spec := style.SpecPaddingRight(size)
				fragments = append(fragments, core.NewFragment(col[y]).AddSpec(spec))
			} else {
				spec := style.SpecRepeatRight(size)
				fragments = append(fragments, core.NewFragment("").AddSpec(spec))
			}

			if i < len(headers)-1 {
				fragments = append(fragments, core.NewFragment(separator.Center))
			}
		}

		fragments = append(fragments, core.NewFragment(separator.Right))

		lines[y] = core.LineFromFragments(fragments...)
	}

	return lines
}

func renderedRowSize(size map[string]int, separator table.SeparatorMeta) int {
	joinSize := (len(size) - 1) * len(separator.Center)
	borderSize := len(separator.Right) + len(separator.Right)

	total := 0
	for _, v := range size {
		total += v
	}

	return total + joinSize + borderSize
}

func adjustSize(size map[string]int, cols int, rowSize int) (map[string]int, bool) {
	if rowSize <= cols {
		return size, true
	}

	excess := rowSize - cols

	h := structure.NewMaxHeapBy(func(c col) int {
		return c.size
	})

	for k, v := range size {
		h.Push(col{k, v})
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

func splitTable(size map[string]int, cols int) []map[string]int {
	tables := make([]map[string]int, 0)

	table := make(map[string]int)
	count := 0

	for k := range size {
		v := min(size[k], cols)

		if count+v >= cols && len(table) > 0 {
			tables = append(tables, table)

			table = make(map[string]int)
			count = 0
		}

		table[k] = v
		count += v
	}

	if len(table) != 0 {
		tables = append(tables, table)
	}

	return tables
}
