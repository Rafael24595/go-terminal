package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/modal"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

const default_modal_menu_name = "IndexMenu"

var modal_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionAll)...,
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
	return modal_definition
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
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *ModalMenu) view(stt state.UIState) viewmodel.ViewModel {
	vm := viewmodel.ViewModelFromUIState(stt)

	frags := input.FragmentFromMenuOption(c.options...)

	modal := modal.NewModalDrawable().
		AddText(c.text...).
		AddOptions(frags...).
		DefineCursor(c.cursor).
		ToDrawable()

	vm.Kernel.Push(modal)

	return *vm
}
