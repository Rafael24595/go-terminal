package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/modal"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

var modal_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionAll)...,
)

type ModalOption struct {
	Fragment text.Fragment
	Action   func() screen.Screen
}

type ModalMenu struct {
	reference string
	text      []text.Line
	options   []ModalOption
	cursor    uint
}

func NewModalMenu() *ModalMenu {
	return &ModalMenu{
		reference: default_index_menu_name,
		text:      make([]text.Line, 0),
		options:   make([]ModalOption, 0),
		cursor:    0,
	}
}

func fragmentFromModalMenu(options ...ModalOption) []text.Fragment {
	frags := make([]text.Fragment, len(options))
	for i := range options {
		frags[i] = options[i].Fragment
	}
	return frags
}

func (c *ModalMenu) SetName(name string) *ModalMenu {
	c.reference = name
	return c
}

func (c *ModalMenu) AddText(text ...text.Line) *ModalMenu {
	c.text = append(c.text, text...)
	return c
}

func (c *ModalMenu) AddOptions(options ...ModalOption) *ModalMenu {
	c.options = append(c.options, options...)
	return c
}

func (c *ModalMenu) SetCursor(cursor uint) *ModalMenu {
	maxIdx := math.SubClampZero(len(c.options), 1)
	c.cursor = math.Clamp(cursor, uint(0), uint(maxIdx))
	return c
}

func (c *ModalMenu) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *ModalMenu) definition() screen.Definition {
	return modal_definition
}

func (c *ModalMenu) name() string {
	return c.reference
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

func (c *ModalMenu) view(stt state.UIState) core.ViewModel {
	vm := core.ViewModelFromUIState(stt)

	frags := fragmentFromModalMenu(c.options...)

	modal := modal.NewModalDrawable().
		AddText(c.text...).
		AddOptions(frags...).
		DefineCursor(c.cursor).
		ToDrawable()

	vm.Lines.Shift(modal)

	return *vm
}
