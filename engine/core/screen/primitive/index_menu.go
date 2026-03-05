package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/indexmenu"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

const default_index_menu_name = "IndexMenu"

var index_menu_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	)...,
)

type MenuOption struct {
	text   text.Fragment
	action func() screen.Screen
}

func NewMenuOption(option text.Fragment, action func() screen.Screen) MenuOption {
	return MenuOption{
		text:   option,
		action: action,
	}
}

func NewMenuOptions(options ...MenuOption) []MenuOption {
	return options
}

func fragmentFromMenuOption(options ...MenuOption) []text.Fragment {
	lines := make([]text.Fragment, len(options))
	for i := range options {
		lines[i] = options[i].text
	}
	return lines
}

type IndexMenu struct {
	reference string
	index     marker.IndexMeta
	title     []text.Line
	options   []MenuOption
	cursor    uint
}

func NewIndexMenu() *IndexMenu {
	return &IndexMenu{
		reference: default_index_menu_name,
		index:     marker.HyphenIndex,
		title:     make([]text.Line, 0),
		options:   make([]MenuOption, 0),
		cursor:    0,
	}
}

func (c *IndexMenu) SetName(name string) *IndexMenu {
	c.reference = name
	return c
}

func (c *IndexMenu) SetIndex(index marker.IndexMeta) *IndexMenu {
	c.index = index
	return c
}

func (c *IndexMenu) AddTitle(title ...text.Line) *IndexMenu {
	c.title = append(c.title, title...)
	return c
}

func (c *IndexMenu) AddOptions(options ...MenuOption) *IndexMenu {
	c.options = append(c.options, options...)
	return c
}

func (c *IndexMenu) SetCursor(cursor uint) *IndexMenu {
	maxIdx := math.SubClampZero(len(c.options), 1)
	c.cursor = math.Clamp(cursor, uint(0), uint(maxIdx))
	return c
}

func (c *IndexMenu) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *IndexMenu) definition() screen.Definition {
	return index_menu_definition
}

func (c *IndexMenu) name() string {
	return c.reference
}

func (c *IndexMenu) update(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	size := uint(len(c.options))
	if size == 0 {
		return screen.EmptyScreenResult()
	}

	switch event.Key.Code {
	case key.ActionArrowUp:
		c.cursor = (c.cursor + size - 1) % size
	case key.ActionTab, key.ActionArrowDown:
		c.cursor = (c.cursor + 1) % size
	case key.ActionEnter:
		option := c.options[c.cursor]
		if option.action != nil {
			scrn := c.options[c.cursor].action()
			return screen.ScreenResult{
				Screen: &scrn,
			}
		}

		assert.Unreachable(
			"menu actions should not be nil: %s - %s",
			c.reference,
			option.text.Text,
		)
	}

	return screen.EmptyScreenResult()
}

func (c *IndexMenu) view(stt state.UIState) core.ViewModel {
	frags := fragmentFromMenuOption(c.options...)

	indexmenu := indexmenu.NewIndexMenuDrawable(c.index, frags).
		Cursor(c.cursor)

	vm := core.ViewModelFromUIState(stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(c.title...),
	)
	vm.Lines.Shift(
		indexmenu.ToDrawable(),
	)

	vm.SetStrategy(
		state.NewFocusPager(),
	)

	option := min(len(c.options)-1, int(c.cursor))
	text := c.options[option].text.Text
	input := core.NewInputLine(line.EagerDrawableFromString(text))
	vm.SetInput(input)

	return *vm
}
