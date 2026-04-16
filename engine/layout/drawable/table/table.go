package table

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/table"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameTableDrawable = "TableDrawable"

type TableDrawable struct {
	loaded     bool
	lazyLoaded bool
	size       terminal.Winsize
	padding    style.HorizontalPosition
	spec       style.Spec
	table      table.Table
	sections   []section
	cursor     input.MatrixCursor
}

func NewTableDrawable(table table.Table, cursor input.MatrixCursor, padding style.HorizontalPosition) *TableDrawable {
	return &TableDrawable{
		loaded:     false,
		lazyLoaded: false,
		size:       terminal.Winsize{},
		padding:    padding,
		spec:       style.SpecEmpty(),
		table:      table,
		sections:   make([]section, 0),
		cursor:     cursor,
	}
}

func TableDrawableFromTable(table table.Table, cursor input.MatrixCursor, padding style.HorizontalPosition) drawable.Drawable {
	return NewTableDrawable(table, cursor, padding).ToDrawable()
}

func (d *TableDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameTableDrawable,
		Code: "",
		Tags: make(set.Set[string]),
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *TableDrawable) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *TableDrawable) wipe() {
	d.lazyLoaded = false
}

func (d *TableDrawable) lazyInit(size terminal.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	d.size = size

	d.spec = makeSpec(d.spec, size, d.padding)
	d.sections = makeSections(d.table, d.cursor, size)

	for i := range d.sections {
		d.sections[i].header.Init()
		d.sections[i].rows.Init()
		d.sections[i].footer.Init()
	}
}

func (d *TableDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]text.Line, 0), false
	}

	d.lazyInit(size)

	headers, footers, remaining := d.drawStatic()
	bodies, hasNext := d.drawDynamic(remaining)

	result := make([]text.Line, 0)
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

func (d *TableDrawable) fillRest(result []text.Line) []text.Line {
	resultSize := min(int(d.size.Rows), len(result))
	for range int(d.size.Rows) - resultSize {
		result = append(result, *text.NewLine(""))
	}

	return result
}

func (d *TableDrawable) drawStatic() ([][]text.Line, [][]text.Line, int) {
	headers := make([][]text.Line, len(d.sections))
	footers := make([][]text.Line, len(d.sections))

	remaining := int(d.size.Rows)
	for i, s := range d.sections {
		header, _ := s.header.Draw(d.size)
		headers[i] = header

		footer, _ := s.footer.Draw(d.size)
		footers[i] = footer

		remaining -= (len(header) + len(footer))
	}

	return headers, footers, remaining
}

func (d *TableDrawable) drawDynamic(remaining int) ([][]text.Line, bool) {
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

			lines, status := s.rows.Draw(d.size)
			if !status {
				empty[i] = 1
			}

			bodies[i] = append(bodies[i], lines...)

			fixRemaining -= len(lines)
		}
	}

	return bodies, len(empty) != len(d.sections)
}

func makeSpec(base style.Spec, size terminal.Winsize, padding style.HorizontalPosition) style.Spec {
	cols := uint(size.Cols)

	var spec style.Spec
	switch padding {
	case style.Left:
		spec = style.SpecPaddingLeft(cols)
	case style.Center:
		spec = style.SpecPaddingCenter(cols)
	case style.Right:
		spec = style.SpecPaddingRight(cols)
	default:
		return base
	}

	return style.MergeSpec(base, spec)
}

func addStyle(spec style.Spec, lines ...text.Line) []text.Line {
	for i := range lines {
		lines[i].Spec = style.MergeSpec(lines[i].Spec, spec)
	}
	return lines
}
