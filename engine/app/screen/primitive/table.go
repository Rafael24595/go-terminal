package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	drawable_table "github.com/Rafael24595/go-terminal/engine/layout/drawable/table"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/model/table"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

const default_table_name = "Table"

var table_navigation_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(
		key.ActionEsc,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	)...,
)

var table_read_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionEnter)...,
)

type Table[T any] struct {
	reference string
	action    *input.TableAction
	title     []text.Line
	table     *table.Table
	cursor    *input.MatrixCursor
	padding   style.HorizontalPosition
}

func NewTable[T any]() *Table[T] {
	return &Table[T]{
		reference: default_table_name,
		action:    input.NewTableAction(),
		title:     make([]text.Line, 0),
		table:     table.NewTable(),
		cursor:    input.NewMatrixCursor(0, 0, false),
		padding:   style.Right,
	}
}

func (c *Table[T]) SetName(name string) *Table[T] {
	c.reference = name
	return c
}

func (c *Table[T]) EnableAction() *Table[T] {
	c.action.EnableMode = true
	return c
}

func (c *Table[T]) DisableAction() *Table[T] {
	c.action.EnableMode = false
	return c
}

func (c *Table[T]) SetActionHandler(handler input.TableActionHandler) *Table[T] {
	c.action.Handler = handler
	return c
}

func (c *Table[T]) DefinePadding(padding style.HorizontalPosition) *Table[T] {
	c.padding = padding
	return c
}

func (c *Table[T]) AddTitle(title ...text.Line) *Table[T] {
	c.title = append(c.title, title...)
	return c
}

func (c *Table[T]) DefineHeaders(headers ...string) *Table[T] {
	c.table = table.NewTable()
	c.table.SetHeaders(headers...)
	return c
}

func (c *Table[T]) AddItems(parser func(T) []table.Field, items ...T) *Table[T] {
	rows := c.table.Rows()
	for i, item := range items {
		for _, field := range parser(item) {
			c.table.SetCell(field.Header, rows+i, field.Value)
		}
	}
	return c
}

func (c *Table[T]) ToScreen() screen.Screen {
	screen := screen.Screen{
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}

	return screen.SetName(c.reference).
		StackFromName()
}

func (c *Table[T]) definition() screen.Definition {
	if !c.action.EnableMode {
		return screen.DefinitionFromKeys()
	}

	if c.action.ActionMode {
		return table_navigation_definition
	}

	return table_read_definition
}

func (c *Table[T]) update(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	state.Pager.ShowPage = true

	if !c.action.EnableMode {
		return screen.ScreenResultFromUIState(state)
	}

	if !c.action.ActionMode {
		return c.updateRead(state, evnt)
	}
	return c.updateNavigation(state, evnt)
}

func (c *Table[T]) updateNavigation(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEsc:
		c.action.ActionMode = false
		c.cursor.Show = c.action.ActionMode
	case key.ActionArrowLeft:
		c.cursor.DecCol()
	case key.ActionArrowRight:
		c.cursor.IncCol(uint32(c.table.Cols() - 1))
	case key.ActionArrowUp:
		c.cursor.DecRow()
	case key.ActionArrowDown:
		c.cursor.IncRow(uint32(c.table.Rows() - 1))
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *Table[T]) updateRead(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEnter:
		c.action.ActionMode = true
		c.cursor.Show = c.action.ActionMode
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *Table[T]) view(stt state.UIState) viewmodel.ViewModel {
	vm := viewmodel.ViewModelFromUIState(stt)

	vm.Header.Push(
		line.EagerDrawableFromLines(c.title...),
	)
	vm.Kernel.Push(
		drawable_table.TableDrawableFromTable(*c.table, *c.cursor, c.padding),
	)

	var input *viewmodel.InputLine
	strategy := state.NewPagePager()
	if c.action.EnableMode && c.action.ActionMode {
		strategy = state.NewFocusPager()

		cell, _ := c.table.FindCellByCoords(int(c.cursor.Row), int(c.cursor.Col))
		input = viewmodel.NewInputLine(line.EagerDrawableFromString(cell))
	}

	vm.SetInput(input)
	vm.SetStrategy(strategy)

	return *vm
}
