package table

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "table_unit"

type TableUnit struct {
	loaded     bool
	lazyLoaded bool
	size       winsize.Winsize
	table      table.Table
	sections   []section
	cursor     input.MatrixCursor
}

func New(table table.Table, cursor input.MatrixCursor) *TableUnit {
	return &TableUnit{
		loaded:     false,
		lazyLoaded: false,
		size:       winsize.Winsize{},
		table:      table,
		sections:   make([]section, 0),
		cursor:     cursor,
	}
}

func UnitFromTable(table table.Table, cursor input.MatrixCursor) drawable.Unit {
	return New(table, cursor).ToUnit()
}

func (d *TableUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(d.init).
		Wipe(d.wipe).
		Draw(d.draw).
		ToUnit()
}

func (d *TableUnit) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *TableUnit) wipe() {
	d.lazyLoaded = false
}

func (d *TableUnit) lazyInit(size winsize.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	d.size = size
	d.sections = makeSections(d.table, d.cursor, size)

	for i := range d.sections {
		d.sections[i].header.Drawable.Init()
		d.sections[i].rows.Drawable.Init()
		d.sections[i].footer.Drawable.Init()
	}
}

func (d *TableUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]text.Line, 0), false
	}

	d.lazyInit(size)

	headers, footers, remaining := d.drawStatic()
	bodies, hasNext := d.drawDynamic(remaining)

	result := make([]text.Line, size.Rows)
	cursor := 0

	for i, body := range bodies {
		if len(body) == 0 {
			continue
		}

		cursor += copy(result[cursor:], headers[i])
		cursor += copy(result[cursor:], body)
		cursor += copy(result[cursor:], footers[i])
	}

	return result, hasNext
}

func (d *TableUnit) drawStatic() ([][]text.Line, [][]text.Line, int) {
	headers := make([][]text.Line, len(d.sections))
	footers := make([][]text.Line, len(d.sections))

	remaining := int(d.size.Rows)
	for i, s := range d.sections {
		header, _ := s.header.Drawable.Draw(d.size)
		headers[i] = header

		footer, _ := s.footer.Drawable.Draw(d.size)
		footers[i] = footer

		remaining -= (len(header) + len(footer))
	}

	return headers, footers, remaining
}

func (d *TableUnit) drawDynamic(remaining int) ([][]text.Line, bool) {
	empty := make(map[int]int)

	sections := len(d.sections)
	if sections == 0 {
		return make([][]text.Line, 0), false
	}

	fixRemaining := remaining - (remaining % sections)

	bodies := make([][]text.Line, sections)
	for fixRemaining > 0 && len(empty) != sections {
		for i, s := range d.sections {
			if fixRemaining <= 0 {
				break
			}

			if _, exists := empty[i]; exists {
				continue
			}

			lines, status := s.rows.Drawable.Draw(d.size)
			if !status {
				empty[i] = 1
			}

			bodies[i] = append(bodies[i], lines...)

			fixRemaining -= len(lines)
		}
	}

	return bodies, len(empty) != len(d.sections)
}
