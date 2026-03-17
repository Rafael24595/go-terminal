package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/indexmenu"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/platform/assert"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/text"
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

type IndexMenu struct {
	reference string
	meta      marker.IndexMeta
	title     []text.Line
	options   []input.MenuOption
	cursor    uint
}

func NewIndexMenu() *IndexMenu {
	return &IndexMenu{
		reference: default_index_menu_name,
		meta:      marker.HyphenIndex,
		title:     make([]text.Line, 0),
		options:   make([]input.MenuOption, 0),
		cursor:    0,
	}
}

func (c *IndexMenu) SetName(name string) *IndexMenu {
	c.reference = name
	return c
}

func (c *IndexMenu) SetMeta(meta marker.IndexMeta) *IndexMenu {
	c.meta = meta
	return c
}

func (c *IndexMenu) AddTitle(title ...text.Line) *IndexMenu {
	c.title = append(c.title, title...)
	return c
}

func (c *IndexMenu) AddOptions(options ...input.MenuOption) *IndexMenu {
	c.options = append(c.options, options...)
	return c
}

func (c *IndexMenu) SetCursor(cursor uint) *IndexMenu {
	maxIdx := math.SubClampZero(len(c.options), 1)
	c.cursor = math.Clamp(cursor, uint(0), uint(maxIdx))
	return c
}

func (c *IndexMenu) ToScreen() screen.Screen {
	screen := screen.Screen{
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}

	return screen.SetName(c.reference).
		StackFromName()
}

func (c *IndexMenu) definition() screen.Definition {
	return index_menu_definition
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
		if option.Action().Name != nil {
			scrn := c.options[c.cursor].Action()
			return screen.ScreenResult{
				Screen: &scrn,
			}
		}

		assert.Unreachable(
			"menu actions should not be nil: %s - %s",
			c.reference,
			option.Label.Text,
		)
	}

	return screen.EmptyScreenResult()
}

func (c *IndexMenu) view(stt state.UIState) viewmodel.ViewModel {
	frags := input.FragmentFromMenuOption(c.options...)

	indexmenu := indexmenu.NewIndexMenuDrawable(frags).
		Meta(c.meta).
		Cursor(c.cursor)

	vm := viewmodel.ViewModelFromUIState(stt)

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
	text := c.options[option].Label.Text
	input := viewmodel.NewInputLine(line.EagerDrawableFromString(text))
	vm.SetInput(input)

	return *vm
}
