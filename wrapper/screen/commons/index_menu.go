package wrapper_commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	wrapper_terminal "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

const default_index_menu_name = "IndexMenu"

type MenuOption struct {
	line   core.Line
	action func() core.Screen
}

func NewMenuOption(line core.Line, action func() core.Screen) MenuOption {
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

func (c *IndexMenu) ToScreen() core.Screen {
	return core.Screen{
		Name:   c.name,
		Update: c.update,
		View:   c.view,
	}
}

func (c *IndexMenu) name() string {
	return c.reference
}

func (c *IndexMenu) update(state state.UIState, event core.ScreenEvent) core.ScreenResult {
	size := uint(len(c.options))
	if size == 0 {
		return core.ScreenResultFromState(state)
	}

	switch event.Key {
	case wrapper_terminal.ARROW_UP:
		c.cursor = (c.cursor + size - 1) % size
	case wrapper_terminal.TAB, wrapper_terminal.ARROW_DOWN:
		c.cursor = (c.cursor + 1) % size
	case "\n", "\r":
		option := c.options[c.cursor]
		if option.action != nil {
			screen := c.options[c.cursor].action()
			return core.ScreenResult{
				Screen: &screen,
			}
		}

		assert.Unreachablef(
			"menu actions should not be nil: %s - %s",
			c.reference,
			option.line,
		)
	}

	return core.ScreenResultFromState(state)
}

func (c *IndexMenu) view(state state.UIState) core.ViewModel {
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

	return core.ViewModel{
		Cursor: &c.cursor,
		Lines:  lines,
	}
}
