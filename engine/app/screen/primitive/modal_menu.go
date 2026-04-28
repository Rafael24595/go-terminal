package primitive

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/modal"
	"github.com/Rafael24595/go-reacterm-core/engine/model/help"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const default_modal_menu_name = "ModalMenu"

const ArgIdModalMenu param.Typed[string] = "id_modal_menu"

var modal_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Active selected"},
	},
	[]key.KeyAction{
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
	cursor    uint
}

func NewModalMenu() *ModalMenu {
	return &ModalMenu{
		reference: default_modal_menu_name,
		text:      make([]text.Line, 0),
		options:   make([]input.MenuOption, 0),
		cursor:    0,
	}
}

func (c *ModalMenu) SetName(name string) *ModalMenu {
	c.reference = name
	return c
}

func (c *ModalMenu) AddText(text ...text.Line) *ModalMenu {
	c.text = append(c.text, text...)
	return c
}

func (c *ModalMenu) AddOptions(options ...input.MenuOption) *ModalMenu {
	c.options = append(c.options, options...)
	return c
}

func (c *ModalMenu) SetCursor(cursor uint) *ModalMenu {
	maxIdx := math.SubClampZero(len(c.options), 1)
	c.cursor = math.Clamp(cursor, uint(0), uint(maxIdx))
	return c
}

func (c *ModalMenu) ToScreen() screen.Screen {
	screen := screen.Screen{
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}

	return screen.SetName(c.reference).
		StackFromName()
}

func (c *ModalMenu) definition() screen.Definition {
	return modal_definition.Definition
}

func (c *ModalMenu) update(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionArrowUp:
		c.cursor = 0
	case key.ActionArrowDown:
		c.cursor = uint(max(0, len(c.options)-1))
	case key.ActionArrowLeft:
		c.cursor = math.SubClampZero(c.cursor, 1)
	case key.ActionArrowRight:
		size := uint(len(c.options))
		if size > 0 {
			c.cursor = min(size-1, c.cursor+1)
		}
	case key.ActionEnter:
		return c.actionEnter(state)
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *ModalMenu) actionEnter(stt *state.UIState) screen.ScreenResult {
	option := c.options[c.cursor]

	state.PushParam(
		stt.Stack,
		c.reference,
		ArgIdModalMenu,
		option.Id,
	)

	if option.Action().Name != nil {
		scrn := c.options[c.cursor].Action()
		return screen.ScreenResultFromScreen(&scrn)
	}

	assert.Unreachable(
		"menu actions should not be nil: %s - %s",
		c.reference,
		option.Label.Text,
	)

	return screen.EmptyScreenResult()
}

func (c *ModalMenu) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	frags := input.FragmentFromMenuOption(c.options...)

	modal := modal.NewModalDrawable().
		AddText(c.text...).
		AddOptions(frags...).
		DefineCursor(c.cursor).
		ToDrawable()

	vm.Kernel.Push(modal)

	vm.Helper.Push(
		key.ActionsToHelpWithOverride(
			modal_definition.Overrides, modal_definition.Actions...,
		)...,
	)

	return *vm
}
