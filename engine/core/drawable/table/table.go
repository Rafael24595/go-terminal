package table

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type TablePadding uint

const (
	Left TablePadding = iota
	Center
	Right
)

type TableDrawable struct {
	initialized bool
	padding     TablePadding
	spec        style.Spec
	table       table.Table
	size        terminal.Winsize
	sections    []section
	cursor      Cursor
}

func NewTableDrawable(table table.Table, cursor Cursor, padding TablePadding) *TableDrawable {
	return &TableDrawable{
		initialized: false,
		padding:     padding,
		spec:        style.SpecEmpty(),
		table:       table,
		size:        terminal.Winsize{},
		sections:    make([]section, 0),
		cursor:      cursor,
	}
}

func TableDrawableFromTable(table table.Table, cursor Cursor, padding TablePadding) core.Drawable {
	return NewTableDrawable(table, cursor, padding).ToDrawable()
}

func (d *TableDrawable) init(size terminal.Winsize) {
	d.initialized = true

	d.spec = makeSpec(d.spec, size, d.padding)
	d.size = size
	d.sections = makeSections(d.table, d.cursor, size)

	for i := range d.sections {
		d.sections[i].header.Init(size)
		d.sections[i].rows.Init(size)
		d.sections[i].footer.Init(size)
	}
}

func (d *TableDrawable) draw() ([]core.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	headers, footers, remaining := d.drawStatic()
	bodies, hasNext := d.drawDynamic(remaining)

	result := make([]core.Line, 0)
	for i, body := range bodies {
		if len(body) == 0 {
			continue
		}

		formatHeaders := addStyle(d.spec, headers[i]...)
		formatBody := addStyle(d.spec, body...)
		formatFooter := addStyle(d.spec, footers[i]...)

		result = append(result, formatHeaders...)
		result = append(result, formatBody...)
		result = append(result, formatFooter...)
	}

	result = d.fillRest(result)
	return result, hasNext
}

func (d *TableDrawable) fillRest(result []core.Line) []core.Line {
	resultSize := min(int(d.size.Rows), len(result))
	for range int(d.size.Rows) - resultSize {
		result = append(result, core.LineFromString(""))
	}

	return result
}

func (d *TableDrawable) drawStatic() ([][]core.Line, [][]core.Line, int) {
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

	return headers, footers, remaining
}

func (d *TableDrawable) drawDynamic(remaining int) ([][]core.Line, bool) {
	empty := make(map[int]int)

	fixRemaining := remaining - (remaining % len(d.sections))

	bodies := make([][]core.Line, len(d.sections))
	for fixRemaining > 0 && len(empty) != len(d.sections) {
		for i, s := range d.sections {
			if fixRemaining <= 0 {
				break
			}

			if _, exists := empty[i]; exists {
				continue
			}

			lines, status := s.rows.Draw()
			if !status {
				empty[i] = 1
			}

			bodies[i] = append(bodies[i], lines...)

			fixRemaining -= len(lines)
		}
	}

	return bodies, len(empty) != len(d.sections)
}

func (d *TableDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}

func makeSpec(base style.Spec, size terminal.Winsize, padding TablePadding) style.Spec {
	cols := uint(size.Cols)

	var spec style.Spec
	switch padding {
	case Left:
		spec = style.SpecPaddingLeft(cols)
	case Center:
		spec = style.SpecPaddingCenter(cols)
	case Right:
		spec = style.SpecPaddingRight(cols)
	default:
		return base
	}

	return style.MergeSpec(base, spec)
}

func addStyle(spec style.Spec, lines ...core.Line) []core.Line {
	for i := range lines {
		lines[i].Spec = style.MergeSpec(lines[i].Spec, spec)
	}
	return lines
}
