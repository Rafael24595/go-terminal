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

func (u *TableUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *TableUnit) init() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *TableUnit) wipe() {
	u.lazyLoaded = false
}

func (u *TableUnit) lazyInit(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	u.size = size
	u.sections = makeSections(u.table, u.cursor, size)

	for i := range u.sections {
		u.sections[i].header.Drawable.Init()
		u.sections[i].rows.Drawable.Init()
		u.sections[i].footer.Drawable.Init()
	}
}

func (u *TableUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]text.Line, 0), false
	}

	u.lazyInit(size)

	headers, footers, remaining := u.drawStatic()
	bodies, hasNext := u.drawDynamic(remaining)

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

func (u *TableUnit) drawStatic() ([][]text.Line, [][]text.Line, int) {
	headers := make([][]text.Line, len(u.sections))
	footers := make([][]text.Line, len(u.sections))

	remaining := int(u.size.Rows)
	for i, s := range u.sections {
		header, _ := s.header.Drawable.Draw(u.size)
		headers[i] = header

		footer, _ := s.footer.Drawable.Draw(u.size)
		footers[i] = footer

		remaining -= (len(header) + len(footer))
	}

	return headers, footers, remaining
}

func (u *TableUnit) drawDynamic(remaining int) ([][]text.Line, bool) {
	empty := make(map[int]int)

	sections := len(u.sections)
	if sections == 0 {
		return make([][]text.Line, 0), false
	}

	fixRemaining := remaining - (remaining % sections)

	bodies := make([][]text.Line, sections)
	for fixRemaining > 0 && len(empty) != sections {
		for i, s := range u.sections {
			if fixRemaining <= 0 {
				break
			}

			if _, exists := empty[i]; exists {
				continue
			}

			lines, status := s.rows.Drawable.Draw(u.size)
			if !status {
				empty[i] = 1
			}

			bodies[i] = append(bodies[i], lines...)

			fixRemaining -= len(lines)
		}
	}

	return bodies, len(empty) != len(u.sections)
}
