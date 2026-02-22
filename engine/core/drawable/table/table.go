package table

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type TableDrawable struct {
	initialized bool
	table       table.Table
	size        terminal.Winsize
	sections    []section
}

func NewTableDrawable(table table.Table) *TableDrawable {
	return &TableDrawable{
		initialized: false,
		table:       table,
		size:        terminal.Winsize{},
		sections:    make([]section, 0),
	}
}

func TableDrawableFromTable(table table.Table) core.Drawable {
	return NewTableDrawable(table).ToDrawable()
}

func (d *TableDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.size = size
	d.sections = makeSections(d.table, size)

	for i := range d.sections {
		d.sections[i].header.Init(size)
		d.sections[i].rows.Init(size)
		d.sections[i].footer.Init(size)
	}
}

func (d *TableDrawable) draw() ([]core.Line, bool) {
	assert.AssertTrue(d.initialized, "the drawable should be initialized before draw")

	headers := make([][]core.Line, len(d.sections))
	footers := make([][]core.Line, len(d.sections))

	remaining := int(d.size.Rows)
	for i, s := range d.sections {
		header, _ := s.header.Draw()
		headers[i] = header

		footer, _ := s.footer.Draw()
		footers[i] = footer

		remaining -= (len(header) + len(footer))
	}

	empty := make(map[int]int)

	bodies := make([][]core.Line, len(d.sections))
	for remaining > 0 && len(empty) != len(d.sections) {
		for i, s := range d.sections {
			if _, exists := empty[i]; exists {
				continue
			}

			lines, status := s.rows.Draw()
			if !status {
				empty[i] = 1
			}

			bodies[i] = append(bodies[i], lines...)

			remaining -= len(lines)
		}
	}

	result := make([]core.Line, 0)
	for i, v := range bodies {
		if len(v) == 0 {
			continue
		}

		result = append(result, headers[i]...)
		result = append(result, v...)
		result = append(result, footers[i]...)
	}

	return result, len(result) != 0
}

func (d *TableDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
