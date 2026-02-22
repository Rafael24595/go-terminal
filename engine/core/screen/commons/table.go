package commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	drawable_table "github.com/Rafael24595/go-terminal/engine/core/drawable/table"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/table"
)

const default_table_name = "Table"

type Field struct {
	Header string
	Value  any
}

type Table[T any] struct {
	reference string
	title     []core.Line
	table     *table.Table
}

func NewTable[T any]() *Table[T] {
	return &Table[T]{
		reference: default_table_name,
		title:     make([]core.Line, 0),
		table:     table.NewTable(),
	}
}

func (c *Table[T]) SetName(name string) *Table[T] {
	c.reference = name
	return c
}

func (c *Table[T]) AddTitle(title ...core.Line) *Table[T] {
	c.title = append(c.title, title...)
	return c
}

func (c *Table[T]) DefineHeaders(headers ...string) *Table[T] {
	c.table = table.NewTable()
	c.table.SetHeaders(headers...)
	return c
}

func (c *Table[T]) AddItems(parser func(T) []Field, items ...T) *Table[T] {
	cols := c.table.Cols()
	for i, item := range items {
		for _, field := range parser(item) {
			c.table.Field(field.Header, cols+i, field.Value)
		}
	}
	return c
}

func (c *Table[T]) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *Table[T]) name() string {
	return c.reference
}

func (c *Table[T]) definition() screen.Definition {
	return screen.Definition{}
}

func (c *Table[T]) update(state state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
	return screen.ScreenResultFromUIState(state)
}

func (c *Table[T]) view(state state.UIState) core.ViewModel {
	vm := core.ViewModelFromUIState(state)

	vm.Header.Shift(
		line.EagerDrawableFromLines(c.title...),
	)
	vm.Lines.Shift(
		drawable_table.TableDrawableFromTable(*c.table),
	)

	return *vm
}
