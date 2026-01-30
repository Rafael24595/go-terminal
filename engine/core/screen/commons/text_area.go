package commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
)

const default_text_area_name = "TextArea"

var text_area_definition = screen.Definition{
	RequireKeys: key.NewKeysCode(key.KeyAll),
}

type TextArea struct {
	reference   string
	title       []core.Line
	selectStart uint
	selectEnd   uint
	buffer      []rune
}

func NewTextArea() *TextArea {
	return &TextArea{
		reference:   default_text_area_name,
		title:       make([]core.Line, 0),
		selectStart: 0,
		selectEnd:   0,
		buffer:      make([]rune, 0),
	}
}

func (c *TextArea) SetName(name string) *TextArea {
	c.reference = name
	return c
}

func (c *TextArea) AddTitle(title ...core.Line) *TextArea {
	c.title = append(c.title, title...)
	return c
}

func (c *TextArea) AddText(text string) *TextArea {
	c.buffer = append(c.buffer, []rune(text)...)
	c.selectStart = uint(len(c.buffer))
	c.selectEnd = uint(len(c.buffer))
	return c
}

func (c *TextArea) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *TextArea) name() string {
	return c.reference
}

func (c *TextArea) definition() screen.Definition {
	return text_area_definition
}

func (c *TextArea) update(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	switch event.Key.Code {
	case key.KeyArrowLeft:
		return c.moveLeft(state, event)
	case key.KeyArrowRight:
		return c.moveRight(state, event)
	case key.KeyBackspace:
		return c.deleteSelection(state)
	}

	text := []rune{event.Key.Rune}
	c.buffer = runes.AppendRange(c.buffer, text, c.selectStart, c.selectEnd)

	c.selectStart = c.selectStart + uint(len(text))
	c.selectEnd = c.selectEnd + uint(len(text))

	return screen.ScreenResultFromState(state)
}

func (c *TextArea) moveLeft(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	if !event.Key.Mod.Has(key.ModCtrl) && event.Key.Mod.Has(key.ModShift) {
		c.selectStart = math.SubClampZero(c.selectStart, 1)
		return screen.ScreenResultFromState(state)
	}

	if event.Key.Mod.Has(key.ModCtrl) {
		c.selectStart = runes.LastIndexOf(c.buffer, ' ', c.selectStart)

		if !event.Key.Mod.Has(key.ModShift) {
			c.selectEnd = c.selectStart
		}

		return screen.ScreenResultFromState(state)
	}

	position := math.SubClampZero(c.selectEnd, 1)

	c.selectStart = position
	c.selectEnd = position

	return screen.ScreenResultFromState(state)
}

func (c *TextArea) moveRight(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	if !event.Key.Mod.Has(key.ModCtrl) && event.Key.Mod.Has(key.ModShift) {
		c.selectEnd = math.Clamp(c.selectEnd+1, 0, uint(len(c.buffer)))
		return screen.ScreenResultFromState(state)
	}

	if event.Key.Mod.Has(key.ModCtrl) {
		c.selectEnd = runes.IndexOf(c.buffer, ' ', c.selectEnd)

		if !event.Key.Mod.Has(key.ModShift) {
			c.selectStart = c.selectEnd
		}

		return screen.ScreenResultFromState(state)
	}

	position := math.Clamp(c.selectEnd+1, 0, uint(len(c.buffer)))

	c.selectStart = position
	c.selectEnd = position

	return screen.ScreenResultFromState(state)
}

func (c *TextArea) deleteSelection(state state.UIState) screen.ScreenResult {
	if len(c.buffer) == 0 {
		return screen.ScreenResultFromState(state)
	}

	start := math.SubClampZero(c.selectStart, 1)
	end := c.selectEnd

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)
	c.selectStart = start
	c.selectEnd = start

	return screen.ScreenResultFromState(state)
}

func (c *TextArea) view(state state.UIState) core.ViewModel {
	renderBuffer := c.buffer

	start := math.SubClampZero(c.selectStart, 1)
	end := c.selectEnd

	if len(renderBuffer) == 0 {
		renderBuffer = append(renderBuffer, ' ')
		start = 0
		end = 1
	}

	if len(renderBuffer) == 0 {
		renderBuffer = append(renderBuffer, ' ')
		start = 0
		end = 1
	}

	text := core.FragmentLine(core.ModePadding(core.Right),
		core.NewFragment(
			string(renderBuffer[0:start]),
			core.Join,
		),
		core.NewFragment(
			string(renderBuffer[start:end]),
			core.Select,
			core.Join,
		),
		core.NewFragment(
			string(renderBuffer[end:]),
			core.Join,
		),
	)

	return core.ViewModel{
		Lines: append(c.title, text),
	}
}
