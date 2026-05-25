package modalmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/modal"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "modal_menu"

const ArgActiveOption param.Typed[string] = "id_modal_menu"

var definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Active selected"},
	},
	[]key.Action{
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
		key.CustomActionBack,
	},
)

type ModalMenu struct {
	reference string
	text      []text.Line
	options   []input.MenuOption
	cursor    uint16
}

func New() *ModalMenu {
	return &ModalMenu{
		reference: Name,
		text:      make([]text.Line, 0),
		options:   make([]input.MenuOption, 0),
		cursor:    0,
	}
}

func (n *ModalMenu) SetName(name string) *ModalMenu {
	n.reference = name
	return n
}

func (n *ModalMenu) AddText(text ...text.Line) *ModalMenu {
	n.text = append(n.text, text...)
	return n
}

func (n *ModalMenu) AddOptions(options ...input.MenuOption) *ModalMenu {
	n.options = append(n.options, options...)
	return n
}

func (n *ModalMenu) SetCursor(cursor uint16) *ModalMenu {
	maxIdx := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	n.cursor = math.Clamp(cursor, 0, maxIdx)
	return n
}

func (n *ModalMenu) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Definition(n.definition).
		Update(n.update).
		View(n.view).
		ToNode()
}

func (n *ModalMenu) definition() screen.Definition {
	return definition
}

func (n *ModalMenu) update(state *state.UIState, evnt screen.Event) screen.Result {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionArrowUp:
		n.cursor = 0
	case key.ActionArrowDown:
		n.cursor = math.SubClampZeroAs[int, uint16](len(n.options), 1)
	case key.ActionArrowLeft:
		n.cursor = math.SubClampZero(n.cursor, 1)
	case key.ActionArrowRight:
		last := math.SubClampZeroAs[int, uint16](len(n.options), 1)
		n.cursor = min(last, n.cursor+1)
	case key.ActionEnter:
		return n.actionEnter(state)
	}

	return screen.ResultFromUIState(state)
}

func (n *ModalMenu) actionEnter(stt *state.UIState) screen.Result {
	option := n.options[n.cursor]

	state.PushParam(
		stt.Stack,
		n.reference,
		ArgActiveOption,
		option.Id,
	)

	node := n.options[n.cursor].Action()
	return screen.ResultFromNode(&node)
}

func (n *ModalMenu) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	frags := input.FragmentFromMenuOption(n.options...)

	modal := modal.New().
		AddText(n.text...).
		AddOptions(frags...).
		DefineCursor(n.cursor).
		ToUnit()

	vm.Kernel.Push(modal)

	return *vm
}
