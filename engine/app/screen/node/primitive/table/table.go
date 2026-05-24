package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding/options"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"

	drawable_table "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/table"
)

const Name = "table"

var read_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit mode"},
	},
	[]key.Action{
		key.ActionEnter,
	},
)

var write_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEsc:   {Code: []string{"ESC"}, Detail: "Write Mode"},
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Active selected"},
	},
	[]key.Action{
		key.ActionEsc,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	},
)

type MarshalFunc[T any] func(T) []table.Field

type Table[T any] struct {
	reference string
	action    *input.TableAction
	table     *table.Table
	cursor    *input.MatrixCursor
	positionY style.VerticalPosition
	positionX style.HorizontalPosition
}

func New[T any]() *Table[T] {
	return &Table[T]{
		reference: Name,
		action:    input.NewTableAction(),
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

func (c *Table[T]) SetHeaders(headers ...string) *Table[T] {
	c.table = table.NewTable()
	c.table.SetHeaders(headers...)
	return c
}

func (c *Table[T]) AddItems(marshal MarshalFunc[T], items ...T) *Table[T] {
	rows := c.table.Rows()
	for i, item := range items {
		index := rows + uint16(i)
		for _, field := range marshal(item) {
			c.table.SetCell(field.Header, index, field.Value)
		}
	}
	return c
}

func (c *Table[T]) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.reference).
		NameToStack().
		Definition(c.definition).
		Update(c.update).
		View(c.view).
		ToNode()
}

func (c *Table[T]) definition() screen.Definition {
	if !c.action.EnableMode {
		return screen.EmptyDefinition()
	}

	if c.action.ActionMode {
		return write_definition
	}

	return read_definition
}

func (c *Table[T]) update(stt *state.UIState, evnt screen.Event) screen.Result {
	stt.Pager.ForceShow = true

	if !c.action.EnableMode {
		return screen.ResultFromUIState(stt)
	}

	if !c.action.ActionMode {
		return c.updateRead(stt, evnt)
	}
	return c.updateNavigation(stt, evnt)
}

func (c *Table[T]) updateNavigation(state *state.UIState, evnt screen.Event) screen.Result {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEsc:
		c.action.ActionMode = false
		c.cursor.Show = c.action.ActionMode
	case key.ActionArrowLeft:
		c.cursor.DecCol()
	case key.ActionArrowRight:
		c.cursor.IncCol(
			math.SubClampZero(c.table.Cols(), 1),
		)
	case key.ActionArrowUp:
		c.cursor.DecRow()
	case key.ActionArrowDown:
		c.cursor.IncRow(
			math.SubClampZero(c.table.Rows(), 1),
		)
	}

	return screen.ResultFromUIState(state)
}

func (c *Table[T]) updateRead(state *state.UIState, evnt screen.Event) screen.Result {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEnter:
		c.action.ActionMode = true
		c.cursor.Show = c.action.ActionMode
	}

	return screen.ResultFromUIState(state)
}

func (c *Table[T]) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	table := drawable_table.UnitFromTable(*c.table, *c.cursor)

	position := padding.NewBuilder().
		Y(hint.Maximize[winsize.Rows](), options.WithPosition(c.positionY)).
		X(hint.Maximize[winsize.Cols](), c.positionX).
		ToUnit(table)

	vm.Kernel.Push(position)

	preficate := pager.PredicatePage()
	if c.action.EnableMode && c.action.ActionMode {
		preficate = pager.PredicateFocus()

		cell, _ := c.table.FindCellByCoords(c.cursor.Row, c.cursor.Col)

		vm.Footer.Push(
			inputline.Wrap(
				drain.UnitFromString(cell),
			),
		)
	}

	vm.Pager.SetPredicate(preficate)

	return *vm
}
