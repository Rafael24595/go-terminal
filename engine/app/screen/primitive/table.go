package primitive

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/position"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_table "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/table"
)

const default_table_name = "Table"

var table_disabled_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{},
	[]key.KeyAction{},
)

var table_read_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit mode"},
	},
	[]key.KeyAction{
		key.ActionEnter,
	},
)

var table_write_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionEsc:   {Code: []string{"ESC"}, Detail: "Write Mode"},
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Active selected"},
	},
	[]key.KeyAction{
		key.ActionEsc,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	},
)

type Table[T any] struct {
	reference string
	action    *input.TableAction
	title     []text.Line
	table     *table.Table
	cursor    *input.MatrixCursor
	positionY style.VerticalPosition
	positionX style.HorizontalPosition
}

func NewTable[T any]() *Table[T] {
	return &Table[T]{
		reference: default_table_name,
		action:    input.NewTableAction(),
		title:     make([]text.Line, 0),
		table:     table.NewTable(),
		cursor:    input.NewMatrixCursor(0, 0, false),
		positionY: style.Middle,
		positionX: style.Center,
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

func (c *Table[T]) SetPositionY(position style.VerticalPosition) *Table[T] {
	c.positionY = position
	return c
}

func (c *Table[T]) SetPositionX(position style.HorizontalPosition) *Table[T] {
	c.positionX = position
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

func (c *Table[T]) definitionSource() screen.DefinitionSources {
	if !c.action.EnableMode {
		return table_disabled_definition
	}

	if c.action.ActionMode {
		return table_write_definition
	}

	return table_read_definition
}

func (c *Table[T]) definition() screen.Definition {
	return c.definitionSource().Definition
}

func (c *Table[T]) update(stt *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	stt.Pager.ForceShow = true

	if !c.action.EnableMode {
		return screen.ScreenResultFromUIState(stt)
	}

	if !c.action.ActionMode {
		return c.updateRead(stt, evnt)
	}
	return c.updateNavigation(stt, evnt)
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

func (c *Table[T]) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	source := c.definitionSource()

	vm.Header.Push(
		block.BlockDrawableFromLines(c.title...),
	)

	table := drawable_table.TableDrawableFromTable(*c.table, *c.cursor)

	vm.Kernel.Push(
		position.NewPositionDrawable(table).
			PositionY(c.positionY).
			PositionX(c.positionX).
			ToDrawable(),
	)

	var input *viewmodel.InputLine
	preficate := pager.PredicatePage()
	if c.action.EnableMode && c.action.ActionMode {
		preficate = pager.PredicateFocus()

		cell, _ := c.table.FindCellByCoords(int(c.cursor.Row), int(c.cursor.Col))
		input = viewmodel.NewInputLine(block.BlockDrawableFromString(cell))
	}

	vm.SetInput(input)

	vm.Pager.SetPredicate(preficate)

	vm.Helper.Push(
		key.ActionsToHelpWithOverride(
			source.Overrides, source.Actions...,
		)...,
	)

	return *vm
}
