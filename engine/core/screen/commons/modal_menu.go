package commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

var modal_menu_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionEnter)...,
)

type ModalOption struct {
	fragment core.Fragment
	action   func() screen.Screen
}

type ModalMenu struct {
	reference string
	text      []core.Line
	options   []ModalOption
	cursor    uint
}

func NewModalMenu() *ModalMenu {
	return &ModalMenu{
		reference: default_index_menu_name,
		text:      make([]core.Line, 0),
		options:   make([]ModalOption, 0),
		cursor:    0,
	}
}

func (c *ModalMenu) SetName(name string) *ModalMenu {
	c.reference = name
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
	return modal_menu_definition
}

func (c *ModalMenu) name() string {
	return c.reference
}

func (c *ModalMenu) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	//TODO: ...
	return screen.EmptyScreenResult()
}

func (c *ModalMenu) view(stt state.UIState) core.ViewModel {
	vm := core.ViewModelFromUIState(stt)
	//TODO: ...
	return *vm
}
