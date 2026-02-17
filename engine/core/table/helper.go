package table

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure"
)

//TODO: Use as a argument.
const min_width = 4

type col struct {
	name string
	size int
}

func renderedSize(size map[string]int, separator, leftBorder, rightBorder string) int {
	joinSize := (len(size) - 1) * len(separator)
	borderSize := len(leftBorder) + len(rightBorder)

	total := 0
	for _, v := range size {
		total += v
	}

	return total + joinSize + borderSize
}

func adjustSize(size map[string]int, termWidth int, renderedSize int) map[string]int {
	if renderedSize <= termWidth {
		return size
	}

	excess := renderedSize - termWidth

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

	return newSize
}

func splitTable(size map[string]int, termWidth int) []map[string]int {
	tables := make([]map[string]int, 0)

	table := make(map[string]int)
	count := 0

	for k := range size {
		v := min(size[k], termWidth)

		if count + v >= termWidth && len(table) > 0 {
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
