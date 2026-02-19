package commons

import (
	"strconv"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

const default_index_menu_name = "IndexMenu"

type IndexKind int

const (
	Numeric IndexKind = iota
	Alphabetic
	Custom
)

var NumericIndex = KindIndex(Numeric)
var AlphabeticIndex = KindIndex(Alphabetic)

var GreaterIndex = CustomIndex(">", "-")
var HyphenIndex = CustomIndex("-", ">")

type IndexMeta struct {
	kind   IndexKind
	index  string
	cursor string
}

func KindIndex(kind IndexKind) IndexMeta {
	return IndexMeta{
		kind: kind,
	}
}

func CustomIndex(index string, cursor string) IndexMeta {
	return IndexMeta{
		kind:   Custom,
		index:  index,
		cursor: cursor,
	}
}

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
	index     IndexMeta
	title     []core.Line
	options   []MenuOption
	cursor    uint
}

func NewIndexMenu() *IndexMenu {
	return &IndexMenu{
		reference: default_index_menu_name,
		index:     HyphenIndex,
		title:     make([]core.Line, 0),
		options:   make([]MenuOption, 0),
		cursor:    0,
	}
}

func (c *IndexMenu) SetName(name string) *IndexMenu {
	c.reference = name
	return c
}

func (c *IndexMenu) SetIndex(index IndexMeta) *IndexMenu {
	c.index = index
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

	digits := math.Digits(len(c.options))

	cursor := 0
	found := false
	for i, o := range c.options {
		selector := []core.Fragment{
			c.makeIndex(i, int(digits)),
			core.NewFragment(" "),
		}

		styledLine := core.FragmentLine(
			style.SpecRepeatLeft(2),
			append(selector, o.line.Text...)...,
		)

		lines = append(lines, styledLine)

		if !found {
			cursor += styledLine.Len()
		}

		if i == int(c.cursor) {
			found = true
		}
	}

	vm := core.ViewModelFromUIState(stt)

	vm.Header.Shift(
		line.LinesEagerDrawableFromLines(c.title...),
	)

	vm.Lines.Shift(
		line.LinesEagerDrawableFromLines(lines...),
	)

	vm.SetCursor(
		state.NewCursorState(c.cursor),
	)

	vm.Pager.Enabled = false
	vm.SetCursor(state.NewCursorState(uint(cursor)))

	return *vm
}

func (c *IndexMenu) makeIndex(cursor, digits int) core.Fragment {
	if c.index.kind == Numeric {
		text := helper.Right(strconv.Itoa(cursor+1), digits)
		index := core.NewFragment(text + ".- ")
		if cursor == int(c.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	if c.index.kind == Alphabetic {
		text := helper.Right(helper.NumberToAlpha(cursor), digits)
		index := core.NewFragment(text + ".- ")
		if cursor == int(c.cursor) {
			index.Atom |= style.AtmBold
		}
		return index
	}

	index := c.index.index
	if cursor == int(c.cursor) {
		index = c.index.cursor
	}
	return core.NewFragment(index)
}
