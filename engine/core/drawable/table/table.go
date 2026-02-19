package table

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type TableDrawable struct {
	table table.Table
	layer *core.LayerStack
}

func NewTableDrawable(table table.Table) *TableDrawable {
	return &TableDrawable{
		table: table,
		layer: core.NewLayerStack(),
	}
}

func TableDrawableFromTable(table table.Table) core.Drawable {
	return NewTableDrawable(table).ToDrawable()
}

func (d *TableDrawable) init(size terminal.Winsize) {
	drawables := makeDrawables(d.table, size)
	d.layer = core.NewLayerStack().Shift(drawables...).Init(size)
}

func (d *TableDrawable) draw() ([]core.Line, bool) {
	lines, ok := d.layer.Draw()
	return lines, ok
}

func (d *TableDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: d.init,
		Draw: d.draw,
	}
}
