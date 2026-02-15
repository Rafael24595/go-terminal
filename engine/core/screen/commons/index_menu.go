package commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

const default_index_menu_name = "IndexMenu"

type MenuOption struct {
	line   core.Line
	action func() screen.Screen
}

func NewMenuOption(line core.Line, action func() screen.Screen) MenuOption {
	return MenuOption{
		line:   line,
		action: action,
	}
}

func NewMenuOptions(options ...MenuOption) []MenuOption {
	return options
}

type IndexMenu struct {
	reference string
	title     []core.Line
	options   []MenuOption
	cursor    uint
}

func NewIndexMenu() *IndexMenu {
	return &IndexMenu{
		reference: default_index_menu_name,
		title:     make([]core.Line, 0),
		options:   make([]MenuOption, 0),
		cursor:    0,
	}
}

func (c *IndexMenu) SetName(name string) *IndexMenu {
	c.reference = name
	return c
}

func (c *IndexMenu) AddTitle(title ...core.Line) *IndexMenu {
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
	return screen.Definition{}
}

func (c *IndexMenu) name() string {
	return c.reference
}

func (c *IndexMenu) update(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
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
			option.line,
		)
	}

	return screen.EmptyScreenResult()
}

func (c *IndexMenu) view(stt state.UIState) core.ViewModel {
	lines := make([]core.Line, 0)

	lines = append(lines, c.title...)

	for i, o := range c.options {
		selector := "-"
		if c.cursor == uint(i) {
			selector = ">"
		}

		fr := []core.Fragment{core.NewFragment(selector)}

		styledLine := core.FragmentLine(
			core.CustomPadding(2, 0),
			append(fr, o.line.Text...)...,
		)
		lines = append(lines, styledLine)
	}

	vm := core.ViewModelFromUIState(stt)

	vm.Lines.Shift(
		line.LinesEagerDrawableFromLines(lines...),
	)
	
	vm.SetCursor(
		state.NewCursorState(c.cursor),
	)

	return *vm
}
