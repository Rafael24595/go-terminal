package table

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure"
)

// TODO: Use as a argument.
const min_width = 4

type col struct {
	name string
	size int
}

func renderedRowSize(size map[string]int, separator SeparatorMeta) int {
	joinSize := (len(size) - 1) * len(separator.center)
	borderSize := len(separator.right) + len(separator.right)

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
