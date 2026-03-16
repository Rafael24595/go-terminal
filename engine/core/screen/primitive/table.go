package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	drawable_table "github.com/Rafael24595/go-terminal/engine/core/drawable/table"
	"github.com/Rafael24595/go-terminal/engine/core/input"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/engine/core/text"
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

type tableActionHandler = func(input.MatrixCursor)

func defaultTableHandler(_ input.MatrixCursor) {}

type tableAction struct {
	enable  bool
	mode    bool
	handler tableActionHandler
}

func enabledTableAction(handler tableActionHandler) *tableAction {
	return &tableAction{
		enable:  true,
		mode:    false,
		handler: handler,
	}
}

func disabledTableAction() *tableAction {
	return &tableAction{
		enable:  false,
		mode:    false,
		handler: defaultTableHandler,
	}
}

type Table[T any] struct {
	reference string
	action    *tableAction
	title     []text.Line
	table     *table.Table
	cursor    *input.MatrixCursor
	padding   style.VerticalPosition
}

func NewTable[T any]() *Table[T] {
	return &Table[T]{
		reference: default_table_name,
		action:    disabledTableAction(),
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

func (c *Table[T]) EnableAction(handler tableActionHandler) *Table[T] {
	c.action = enabledTableAction(handler)
	return c
}

func (c *Table[T]) DisableAction() *Table[T] {
	c.action = disabledTableAction()
	return c
}

func (c *Table[T]) DefinePadding(padding style.VerticalPosition) *Table[T] {
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
	if !c.action.enable {
		return screen.DefinitionFromKeys()
	}

	if c.action.mode {
		return table_navigation_definition
	}
	
	return table_read_definition
}

func (c *Table[T]) update(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	state.Pager.ShowPage = true

	if !c.action.enable {
		return screen.ScreenResultFromUIState(state)
	}

	if !c.action.mode {
		return c.updateRead(state, evnt)
	}
	return c.updateNavigation(state, evnt)
}

func (c *Table[T]) updateNavigation(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEsc:
		c.action.mode = false
		c.cursor.Show = c.action.mode
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
		c.action.mode = true
		c.cursor.Show = c.action.mode
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *Table[T]) view(stt state.UIState) core.ViewModel {
	vm := core.ViewModelFromUIState(stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(c.title...),
	)
	vm.Lines.Shift(
		drawable_table.TableDrawableFromTable(*c.table, *c.cursor, c.padding),
	)

	var input *core.InputLine
	strategy := state.NewPagePager()
	if c.action.enable && c.action.mode {
		strategy = state.NewFocusPager()

		cell, _ := c.table.FindCellByCoords(int(c.cursor.Row), int(c.cursor.Col))
		input = core.NewInputLine(line.EagerDrawableFromString(cell))
	}

	vm.SetInput(input)
	vm.SetStrategy(strategy)

	return *vm
}
