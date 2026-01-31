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
	case key.KeyHome:
		c.moveCursorTo(0)
		return screen.ScreenResultFromState(state)
	case key.KeyEnd:
		c.moveCursorTo(uint(len(c.buffer)))
		return screen.ScreenResultFromState(state)
	case key.KeyArrowLeft:
		return c.moveBackward(state, event)
	case key.KeyArrowRight:
		return c.moveForward(state, event)
	case key.KeyBackspace, key.KeyDeleteWordBackward:
		return c.deleteBackward(state, event.Key.Code == key.KeyDeleteWordBackward)
	case key.KeyDelete, key.KeyDeleteWordForward:
		return c.deleteForward(state, event.Key.Code == key.KeyDeleteWordForward)
	case key.KeyArrowUp, key.KeyArrowDown:
	}

	text := []rune{event.Key.Rune}
	c.buffer = runes.AppendRange(c.buffer, text, c.selectStart, c.selectEnd)

	position := c.selectStart + uint(len(text))
	c.moveCursorTo(position)

	return screen.ScreenResultFromState(state)
}

func (c *TextArea) moveBackward(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromState(state)

	if !event.Key.Mod.Has(key.ModCtrl) && event.Key.Mod.Has(key.ModShift) {
		start := math.SubClampZero(c.selectStart, 1)
		c.moveSelectTo(start, c.selectEnd)
		return result
	}

	if event.Key.Mod.Has(key.ModCtrl) {
		start := runes.BackwardIndex(c.buffer, ' ', c.selectStart)

		end := c.selectEnd
		if !event.Key.Mod.Has(key.ModShift) {
			end = start
		}

		c.moveSelectTo(start, end)

		return result
	}

	position := math.SubClampZero(c.selectEnd, 1)
	c.moveCursorTo(position)

	return result
}

func (c *TextArea) moveForward(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromState(state)

	if !event.Key.Mod.Has(key.ModCtrl) && event.Key.Mod.Has(key.ModShift) {
		end := min(uint(len(c.buffer)), c.selectEnd+1)
		c.moveSelectTo(c.selectStart, end)
		return result
	}

	if event.Key.Mod.Has(key.ModCtrl) {
		end := runes.ForwardIndex(c.buffer, ' ', c.selectEnd)

		start := c.selectStart
		if !event.Key.Mod.Has(key.ModShift) {
			start = end
		}

		c.moveSelectTo(start, end)

		return result
	}

	position := min(uint(len(c.buffer)), c.selectEnd+1)
	c.moveCursorTo(position)

	return result
}

func (c *TextArea) deleteBackward(state state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromState(state)

	if len(c.buffer) == 0 {
		return result
	}

	var start uint
	if word {
		start = runes.BackwardIndex(c.buffer, ' ', c.selectStart)
	} else {
		start = math.SubClampZero(c.selectStart, 1)
	}

	end := c.selectEnd

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)

	c.moveCursorTo(start)

	return result
}

func (c *TextArea) deleteForward(state state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromState(state)

	if len(c.buffer) == 0 {
		return result
	}

	var end uint
	if word {
		end = runes.ForwardIndex(c.buffer, ' ', c.selectEnd)
	} else {
		end = min(uint(len(c.buffer)), c.selectEnd+1)
	}

	start := c.selectStart

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)

	c.moveCursorTo(start)

	return result
}

func (c *TextArea) moveCursorTo(position uint) {
	min := uint(1)
	len := uint(len(c.buffer))

	if len == 0 {
		min = 0
	}

	c.selectStart = math.Clamp(position, min, len)
	c.selectEnd = c.selectStart
}

func (c *TextArea) moveSelectTo(start, end uint) {
	min := uint(1)
	len := uint(len(c.buffer))

	if len == 0 {
		min = 0
	}

	c.selectStart = math.Clamp(start, min, len)
	c.selectEnd = math.Clamp(end, min, len)
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
