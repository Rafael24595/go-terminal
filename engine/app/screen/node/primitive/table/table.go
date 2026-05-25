package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
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

func (n *Table[T]) SetName(name string) *Table[T] {
	n.reference = name
	return n
}

func (n *Table[T]) EnableAction() *Table[T] {
	n.action.EnableMode = true
	return n
}

func (n *Table[T]) DisableAction() *Table[T] {
	n.action.EnableMode = false
	return n
}

func (n *Table[T]) SetActionHandler(handler input.TableActionHandler) *Table[T] {
	n.action.Handler = handler
	return n
}

func (n *Table[T]) SetPositionY(position style.VerticalPosition) *Table[T] {
	n.positionY = position
	return n
}

func (n *Table[T]) SetPositionX(position style.HorizontalPosition) *Table[T] {
	n.positionX = position
	return n
}

func (n *Table[T]) SetHeaders(headers ...string) *Table[T] {
	n.table = table.NewTable()
	n.table.SetHeaders(headers...)
	return n
}

func (n *Table[T]) AddItems(marshal MarshalFunc[T], items ...T) *Table[T] {
	rows := n.table.Rows()
	for i, item := range items {
		index := rows + uint16(i)
		for _, field := range marshal(item) {
			n.table.SetCell(field.Header, index, field.Value)
		}
	}
	return n
}

func (n *Table[T]) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Definition(n.definition).
		Update(n.update).
		View(n.view).
		ToNode()
}

func (n *Table[T]) definition() screen.Definition {
	if !n.action.EnableMode {
		return screen.EmptyDefinition()
	}

	if n.action.ActionMode {
		return write_definition
	}

	return read_definition
}

func (n *Table[T]) update(stt *state.UIState, evnt screen.Event) screen.Result {
	stt.Pager.ForceShow = true

	if !n.action.EnableMode {
		return screen.ResultFromUIState(stt)
	}

	if !n.action.ActionMode {
		return n.updateRead(stt, evnt)
	}
	return n.updateNavigation(stt, evnt)
}

func (n *Table[T]) updateNavigation(state *state.UIState, evnt screen.Event) screen.Result {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEsc:
		n.action.ActionMode = false
		n.cursor.Show = n.action.ActionMode
	case key.ActionArrowLeft:
		n.cursor.DecCol()
	case key.ActionArrowRight:
		n.cursor.IncCol(
			math.SubClampZero(n.table.Cols(), 1),
		)
	case key.ActionArrowUp:
		n.cursor.DecRow()
	case key.ActionArrowDown:
		n.cursor.IncRow(
			math.SubClampZero(n.table.Rows(), 1),
		)
	}

	return screen.ResultFromUIState(state)
}

func (n *Table[T]) updateRead(state *state.UIState, evnt screen.Event) screen.Result {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEnter:
		n.action.ActionMode = true
		n.cursor.Show = n.action.ActionMode
	}

	return screen.ResultFromUIState(state)
}

func (n *Table[T]) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	table := drawable_table.UnitFromTable(*n.table, *n.cursor)

	position := padding.NewBuilder().
		Y(hint.Maximize[winsize.Rows](), rows.WithPosition(n.positionY)).
		X(hint.Maximize[winsize.Cols](), cols.WithPosition(n.positionX)).
		ToUnit(table)

	vm.Kernel.Push(position)

	preficate := pager.PredicatePage()
	if n.action.EnableMode && n.action.ActionMode {
		preficate = pager.PredicateFocus()

		cell, _ := n.table.FindCellByCoords(n.cursor.Row, n.cursor.Col)

		vm.Footer.Push(
			inputline.Wrap(
				drain.UnitFromString(cell),
			),
		)
	}

	vm.Pager.SetPredicate(preficate)

	return *vm
}
