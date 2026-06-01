package indexmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/indexmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

const Name = "index_menu"

const ArgActiveIndex param.Typed[string] = "id_index_menu"

var index_menu_definition = screen.DefinitionFromActions(
	[]key.Action{
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
		key.CustomActionPointer,
	}...,
)

type IndexMenu struct {
	reference string
	pointer   uint8
	meta      marker.IndexMeta
	options   []input.MenuOption
	cursor    uint16
}

func New() *IndexMenu {
	return &IndexMenu{
		reference: Name,
		pointer:   0,
		meta:      marker.HyphenIndex,
		options:   make([]input.MenuOption, 0),
		cursor:    0,
	}
}

func (n *IndexMenu) SetName(name string) *IndexMenu {
	n.reference = name
	return n
}

func (n *IndexMenu) SetMeta(meta marker.IndexMeta) *IndexMenu {
	n.meta = meta
	return n
}

func (n *IndexMenu) AddOptions(options ...input.MenuOption) *IndexMenu {
	n.options = append(n.options, options...)
	return n
}

func (n *IndexMenu) SetCursor(cursor uint16) *IndexMenu {
	maxIdx := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	n.cursor = math.Clamp(cursor, 0, maxIdx)
	return n
}

func (n *IndexMenu) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *IndexMenu) keys() screen.Definition {
	return index_menu_definition
}

func (n *IndexMenu) tick(uiState *state.UIState, event screen.Event) screen.Result {
	size := uint16(len(n.options))
	if size == 0 {
		return screen.EmptyResult()
	}

	switch event.Key.Code {
	case key.ActionArrowUp:
		n.cursor = (n.cursor + size - 1) % size
		n.tickToStack(uiState)
	case key.ActionTab, key.ActionArrowDown:
		n.cursor = (n.cursor + 1) % size
		n.tickToStack(uiState)
	case key.ActionEnter:
		n.tickToStack(uiState)
		return n.actionEnter()
	case key.CustomActionPointer:
		n.pointer = indexmenu.NextPointer(n.pointer)
	}

	return screen.EmptyResult()
}

func (n *IndexMenu) tickToStack(uiState *state.UIState) {
	if n.cursor >= uint16(len(n.options)) {
		uiState.Stack.RemoveArgument(
			n.reference,
			string(ArgActiveIndex),
		)
		return
	}

	state.PushParam(
		uiState.Stack,
		n.reference,
		ArgActiveIndex,
		n.options[n.cursor].Id,
	)
}

func (n *IndexMenu) actionEnter() screen.Result {
	node := n.options[n.cursor].Action()
	return screen.ResultFromNode(&node)
}

func (n *IndexMenu) view(_ state.UIState) viewmodel.ViewModel {
	frags := input.FragmentFromMenuOption(n.options...)

	pointer := indexmenu.FindPointer(n.pointer)

	indexmenu := indexmenu.New(frags).
		Pointer(pointer).
		Meta(n.meta).
		Cursor(n.cursor)

	vm := viewmodel.New()

	vm.Kernel.Push(
		indexmenu.ToUnit(),
	)

	index := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	option := min(index, n.cursor)
	text := n.options[option].Label.Text

	vm.Footer.Push(
		inputline.FromString(text),
	)

	vm.Pager.SetPredicate(
		pager.PredicateFocus(),
	)

	return *vm
}
