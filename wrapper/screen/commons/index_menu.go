package wrapper_commons

import (
	"github.com/Rafael24595/go-terminal/engine/core"
	wrapper_terminal "github.com/Rafael24595/go-terminal/wrapper/terminal"
)

type MenuOption struct {
	line core.Line
}

func NewMenuOption(line core.Line) MenuOption {
	return MenuOption{
		line: line,
	}
}

func NewMenuOptions(options ...MenuOption) []MenuOption {
	return options
}

type IndexMenu struct {
	title   []core.Line
	options []MenuOption
	cursor  uint
}

func NewIndexMenu() *IndexMenu {
	return &IndexMenu{
		title:   make([]core.Line, 0),
		options: make([]MenuOption, 0),
		cursor:  0,
	}
}

func (c *IndexMenu) SetCursor(cursor uint) *IndexMenu {
	c.cursor = min(uint(len(c.options)), cursor)
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

func (c *IndexMenu) ToScreen() core.Screen {
	return core.Screen{
		Update: c.Update,
		View:   c.View,
	}
}

func (c *IndexMenu) Update(e core.ScreenEvent) {
	size := uint(len(c.options))
	switch e.Key {
	case wrapper_terminal.ARROW_UP:
		c.cursor = ((c.cursor - 1) % (size * 2)) % size
	case wrapper_terminal.TAB, wrapper_terminal.ARROW_DOWN:
		c.cursor = (c.cursor + 1) % size
	}
}

func (c *IndexMenu) View() core.ViewModel {
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
		Lines: lines,
	}
}
