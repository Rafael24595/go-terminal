package commons

import (
	"slices"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/helper/line"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
)

const default_text_area_name = "TextArea"

var next_word_runes = []runes.RuneDefinition{
	{
		Rune: ' ',
		Skip: false,
	},
	{
		Rune: '.',
		Skip: true,
	},
	{
		Rune: ',',
		Skip: true,
	},
	{
		Rune: key.ENTER_LF,
		Skip: true,
	},
}

var next_line_runes = []runes.RuneDefinition{
	{
		Rune: key.ENTER_LF,
		Skip: true,
	},
}

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
	ky := event.Key

	switch ky.Code {
	case key.KeyHome:
		return c.moveHome(state, event)
	case key.KeyEnd:
		return c.moveEnd(state, event)
	case key.KeyArrowLeft:
		return c.moveBackward(state, event)
	case key.KeyArrowRight:
		return c.moveForward(state, event)
	case key.KeyBackspace, key.KeyDeleteWordBackward:
		return c.deleteBackward(state, ky.Code == key.KeyDeleteWordBackward)
	case key.KeyDelete, key.KeyDeleteWordForward:
		return c.deleteForward(state, ky.Code == key.KeyDeleteWordForward)
	case key.KeyEnter:
		ky = *key.NewKeyRune(key.ENTER_LF)
	case key.KeyArrowUp:
		return c.moveUp(state)
	case key.KeyArrowDown:
		return c.moveDown(state)
	}

	text := []rune{ky.Rune}
	c.buffer = runes.AppendRange(c.buffer, text, c.selectStart, c.selectEnd)

	position := c.selectEnd + uint(len(text))
	c.moveCursorTo(position)

	return screen.ScreenResultFromUIState(state)
}

func (c *TextArea) moveHome(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.Has(key.ModCtrl) {
		c.moveCursorTo(0)
		return result
	}

	start := runes.BackwardIndexWithLimit(c.buffer, next_line_runes, c.selectStart)

	end := c.selectEnd
	if !event.Key.Mod.Has(key.ModShift) {
		end = start
	}

	c.moveSelectTo(start, end)

	return result
}

func (c *TextArea) moveEnd(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.Has(key.ModCtrl) {
		c.moveCursorTo(uint(len(c.buffer)))
		return result
	}

	start := runes.ForwardIndexWithLimit(c.buffer, next_line_runes, c.selectStart)

	end := c.selectEnd
	if !event.Key.Mod.Has(key.ModShift) {
		end = start
	}

	c.moveSelectTo(start, end)

	return result
}

func (c *TextArea) moveUp(state state.UIState) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	distance := line.DistanceFromLF(c.buffer, int(c.selectStart))

	prevLineStart := line.FindLineStart(c.buffer, int(c.selectStart))
	if prevLineStart == 0 {
		return result
	}

	targetLineStart := line.FindLineStart(c.buffer, prevLineStart-1)
	position := line.ClampToLine(c.buffer, targetLineStart, distance)

	c.moveCursorTo(uint(position))

	return result
}

func (c *TextArea) moveDown(state state.UIState) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	distance := line.DistanceFromLF(c.buffer, int(c.selectStart))

	nextLineStart := line.FindNextLineStart(c.buffer, int(c.selectStart))
	if nextLineStart == -1 {
		c.moveCursorTo(uint(len(c.buffer)))
		return result
	}

	position := line.ClampToLine(c.buffer, nextLineStart, distance)

	c.moveCursorTo(uint(position))

	return result
}

func (c *TextArea) moveBackward(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if !event.Key.Mod.Has(key.ModCtrl) && event.Key.Mod.Has(key.ModShift) {
		start := math.SubClampZero(c.selectStart, 1)
		c.moveSelectTo(start, c.selectEnd)
		return result
	}

	if event.Key.Mod.Has(key.ModCtrl) {
		start := runes.BackwardIndex(c.buffer, next_word_runes, c.selectStart)

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
	result := screen.ScreenResultFromUIState(state)

	if !event.Key.Mod.Has(key.ModCtrl) && event.Key.Mod.Has(key.ModShift) {
		end := min(uint(len(c.buffer)), c.selectEnd+1)
		c.moveSelectTo(c.selectStart, end)
		return result
	}

	if event.Key.Mod.Has(key.ModCtrl) {
		end := runes.ForwardIndex(c.buffer, next_word_runes, c.selectEnd)

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
	result := screen.ScreenResultFromUIState(state)

	if len(c.buffer) == 0 {
		return result
	}

	var start uint
	if word {
		start = runes.BackwardIndex(c.buffer, next_word_runes, c.selectStart)
	} else {
		start = math.SubClampZero(c.selectStart, 1)
	}

	end := c.selectEnd

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)

	c.moveCursorTo(start)

	return result
}

func (c *TextArea) deleteForward(state state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if len(c.buffer) == 0 {
		return result
	}

	var end uint
	if word {
		end = runes.ForwardIndex(c.buffer, next_word_runes, c.selectEnd)
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

func (c *TextArea) view(stt state.UIState) core.ViewModel {
	renderBuffer := c.buffer

	start := math.SubClampZero(c.selectStart, 1)
	end := c.selectEnd

	if len(renderBuffer) == 0 {
		renderBuffer = append(renderBuffer, ' ')
		start = 0
		end = 1
	}

	text := core.FragmentLine(core.ModePadding(core.Right))

	beforeSelect := string(renderBuffer[0:start])
	text.Text = append(text.Text, core.NewFragment(beforeSelect))

	onSelect := string(renderBuffer[start:end])
	text.Text = append(text.Text, core.NewFragment(onSelect, core.Select))

	afterSelect := string(renderBuffer[end:])
	if len(afterSelect) > 0 {
		text.Text = append(text.Text, core.NewFragment(afterSelect))
	}

	return *core.ViewModelFromUIState(stt).
		AddHeader(c.title...).
		AddLines(c.normalizeLinesEnd(text)...).
		SetPager(state.EmptyPagerState()).
		SetCursor(state.NewCursorState(c.selectStart))
}

func (c *TextArea) normalizeLinesEnd(text core.Line) []core.Line {
	lines := make([]core.Line, 0)

	currentLine := core.FragmentLine(text.Padding)

	for textIndex, f := range text.Text {
		normalized := runes.NormalizeLineEnd(f.Text)

		parts := strings.Split(normalized, "\n")
		if len(parts) == 1 {
			currentLine.Text = append(
				currentLine.Text,
				core.NewFragment(parts[0], f.Styles...),
			)

			continue
		}

		for partIndex, part := range parts {
			if c.isCaretPrintable(text, textIndex, part, partIndex) {
				part += " "
			}

			currentLine.Text = append(
				currentLine.Text,
				core.NewFragment(part, f.Styles...),
			)

			if partIndex < len(parts)-1 {
				lines = append(lines, currentLine)
				currentLine = core.FragmentLine(text.Padding)
			}
		}
	}

	if len(currentLine.Text) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

func (c *TextArea) isCaretPrintable(text core.Line, textIndex int, part string, partIndex int) bool {
	fragment := text.Text[textIndex]

	isCaret := len(part) == 0 && slices.Contains(fragment.Styles, core.Select)
	if !isCaret {
		return false
	}

	atLineStart := partIndex == 0
	if atLineStart {
		return true
	}

	atBufferEnd := textIndex == len(text.Text)-1
	if atBufferEnd {
		return true
	}

	atEmptyLine := text.Text[textIndex+1].Text[0] == key.ENTER_LF
	return atEmptyLine
}
