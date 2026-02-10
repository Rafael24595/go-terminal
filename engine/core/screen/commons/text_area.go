package commons

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/event"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/helper/line"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/helper/text"
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
		Skip: false,
	},
}

var text_area_definition = screen.Definition{
	RequireKeys: key.NewKeysCode(key.ActionAll),
}

type TextArea struct {
	reference string
	history   *event.TextEventService
	title     []core.Line
	footer    []core.Line
	caret     uint
	anchor    uint
	buffer    []rune
	index     bool
}

func NewTextArea() *TextArea {
	return &TextArea{
		reference: default_text_area_name,
		history:   event.NewTextEventService(),
		title:     make([]core.Line, 0),
		footer:    make([]core.Line, 0),
		caret:     0,
		anchor:    0,
		buffer:    make([]rune, 0),
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

func (c *TextArea) AddFooter(footer ...core.Line) *TextArea {
	c.footer = append(c.footer, footer...)
	return c
}

func (c *TextArea) AddText(text string) *TextArea {
	c.buffer = append(c.buffer, []rune(text)...)
	c.caret = uint(len(c.buffer))
	c.anchor = c.caret
	return c
}

func (c *TextArea) ShowIndex() *TextArea {
	c.index = true
	return c
}

func (c *TextArea) HideIndex() *TextArea {
	c.index = false
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

func (c *TextArea) selectStart() uint {
	if c.anchor < c.caret {
		return c.anchor
	}
	return c.caret
}

func (c *TextArea) selectEnd() uint {
	if c.anchor < c.caret {
		return c.caret
	}
	return c.anchor
}

func (c *TextArea) update(state state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionHome:
		return c.moveHome(state, evnt)
	case key.ActionEnd:
		return c.moveEnd(state, evnt)
	case key.ActionArrowLeft:
		return c.moveBackward(state, evnt)
	case key.ActionArrowRight:
		return c.moveForward(state, evnt)
	case key.ActionBackspace, key.ActionDeleteBackward:
		return c.deleteBackward(state, ky.Code == key.ActionDeleteBackward)
	case key.ActionDelete, key.ActionDeleteForward:
		return c.deleteForward(state, ky.Code == key.ActionDeleteForward)
	case key.ActionEnter:
		ky = *key.NewKeyRune(key.ENTER_LF)
	case key.ActionArrowUp:
		return c.moveUp(state, evnt)
	case key.ActionArrowDown:
		return c.moveDown(state, evnt)
	case key.CustomActionUndo, key.CustomActionRedo:
		return c.undoRedo(state, ky)
	}

	end := c.selectEnd()

	start := c.selectStart()
	fixEnd := end
	if start != end {
		start = math.SubClampZero(start, 1)
		fixEnd += 1
	}

	text := text.FullTextTransformer.Apply(ky, start, end, c.buffer)

	c.history.PushEvent(event.Insert, start, fixEnd, string(c.buffer[start:end]), string(text))

	c.buffer = runes.AppendRange(c.buffer, text, start, end)

	position := start + uint(len(text))
	c.moveCaretTo(position)

	return screen.ScreenResultFromUIState(state)
}

func (c *TextArea) undoRedo(state state.UIState, ky key.Key) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	var delta *event.Delta
	switch ky.Code {
	case key.CustomActionUndo:
		delta = c.history.Undo()
	case key.CustomActionRedo:
		delta = c.history.Redo()
	default:
		assert.Unreachable("unsupported key code '%d'", ky.Code)
		delta = c.history.Redo()
	}

	if delta == nil {
		return result
	}

	c.buffer = event.ApplyDelta(c.buffer, delta)
	newTextRunes := []rune(delta.Text)
	c.moveCaretTo(delta.Start + uint(len(newTextRunes)))

	return result
}

func (c *TextArea) moveHome(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasAny(key.ModCtrl) {
		c.moveCaretTo(0)
		return result
	}

	caret := runes.BackwardIndexWithLimit(c.buffer, next_line_runes, c.caret)

	anchor := c.anchor
	if event.Key.Mod.HasNone(key.ModShift) {
		c.moveCaretTo(caret)
		return result
	}

	c.moveSelectTo(caret, anchor)

	return result
}

func (c *TextArea) moveEnd(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasAny(key.ModCtrl) {
		c.moveCaretTo(uint(len(c.buffer)))
		return result
	}

	caret := runes.ForwardIndexWithLimit(c.buffer, next_line_runes, c.caret)

	anchor := c.anchor
	if event.Key.Mod.HasNone(key.ModShift) {
		c.moveCaretTo(caret)
		return result
	}

	c.moveSelectTo(caret, anchor)

	return result
}

func (c *TextArea) moveUp(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	start := c.caret
	distance := line.DistanceFromLF(c.buffer, int(start))

	prevLineStart := line.FindPrevLineStart(c.buffer, int(start))
	if prevLineStart == -1 {
		if event.Key.Mod.HasAny(key.ModShift) {
			c.moveSelectTo(0, c.anchor)
			return result
		}

		c.moveCaretTo(0)
		return result
	}

	position := line.ClampToLine(c.buffer, prevLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		c.moveSelectTo(uint(position), c.anchor)
	} else {
		c.moveCaretTo(uint(position))
	}

	return result
}

func (c *TextArea) moveDown(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	start := c.caret
	distance := line.DistanceFromLF(c.buffer, int(start))

	nextLineStart := line.FindNextLineStart(c.buffer, int(start))
	if nextLineStart == -1 {
		if event.Key.Mod.HasAny(key.ModShift) {
			c.moveSelectTo(uint(len(c.buffer)), c.anchor)
			return result
		}

		c.moveCaretTo(uint(len(c.buffer)))
		return result
	}

	position := line.ClampToLine(c.buffer, nextLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		c.moveSelectTo(uint(position), c.anchor)
	} else {
		c.moveCaretTo(uint(position))
	}

	return result
}

func (c *TextArea) moveBackward(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := math.SubClampZero(c.caret, 1)
		c.moveCaretTo(caret)
		return result
	}

	anchor := c.anchor
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := math.SubClampZero(c.caret, 1)
		c.moveSelectTo(caret, anchor)
		return result
	}

	caret := runes.BackwardIndex(c.buffer, next_word_runes, c.caret)
	if event.Key.Mod.HasNone(key.ModShift) {
		c.moveCaretTo(caret)
		return result
	}

	c.moveSelectTo(caret, anchor)
	return result
}

func (c *TextArea) moveForward(state state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := min(uint(len(c.buffer)), c.caret+1)
		c.moveCaretTo(caret)
		return result
	}

	anchor := c.anchor
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := min(uint(len(c.buffer)), c.caret+1)
		c.moveSelectTo(caret, anchor)
		return result
	}

	caret := runes.ForwardIndex(c.buffer, next_word_runes, c.caret)
	if event.Key.Mod.HasNone(key.ModShift) {
		c.moveCaretTo(caret)
		return result
	}

	c.moveSelectTo(caret, anchor)
	return result
}

func (c *TextArea) deleteBackward(state state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if len(c.buffer) == 0 {
		return result
	}

	start := c.selectStart()

	if word {
		start = runes.BackwardIndex(c.buffer, next_word_runes, start)
	} else {
		start = math.SubClampZero(start, 1)
	}

	end := c.selectEnd()

	c.history.PushEvent(event.DeleteBackward, start, end, string(c.buffer[start:end]), "")

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)

	c.moveCaretTo(start)
	return result
}

func (c *TextArea) deleteForward(state state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if len(c.buffer) == 0 {
		return result
	}

	end := c.selectEnd()

	if word {
		end = runes.ForwardIndex(c.buffer, next_word_runes, end)
	} else {
		end = min(uint(len(c.buffer)), end+1)
	}

	start := math.SubClampZero(c.selectStart(), 1)

	c.history.PushEvent(event.DeleteForward, start, end, string(c.buffer[start:end]), "")

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)

	c.moveCaretTo(start)
	return result
}

func (c *TextArea) moveCaretTo(caret uint) {
	min := uint(1)
	len := uint(len(c.buffer))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = c.caret
}

func (c *TextArea) moveSelectTo(caret, anchor uint) {
	min := uint(1)
	len := uint(len(c.buffer))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = math.Clamp(anchor, min, len)
}

func (c *TextArea) view(stt state.UIState) core.ViewModel {
	renderBuffer := c.buffer

	start := math.SubClampZero(c.selectStart(), 1)
	end := c.selectEnd()

	if len(renderBuffer) == 0 {
		renderBuffer = append(renderBuffer, []rune(PRINTABLE_CARET)...)
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

	lines := c.normalizeLinesEnd(text)
	lines = c.fixEmptyLines(lines)

	return *core.ViewModelFromUIState(stt).
		AddHeader(c.title...).
		AddLines(lines...).
		AddFooter(c.footer...).
		SetPager(state.EmptyPagerState()).
		SetCursor(state.NewCursorState(c.selectEnd()))
}

func (c *TextArea) normalizeLinesEnd(text core.Line) []core.Line {
	lines := make([]core.Line, 0)

	index := uint16(1)

	currentLine := core.FragmentLine(text.Padding)
	if c.index {
		currentLine.SetOrder(index)
	}

	for textIndex, f := range text.Text {
		normalized := runes.NormalizeLineEnd(f.Text)

		parts := strings.Split(normalized, "\n")
		if len(parts) == 1 {
			currentLine.Text = append(
				currentLine.Text,
				core.NewFragment(parts[0], f.Styles),
			)

			continue
		}

		for partIndex, part := range parts {
			if c.isCaretPrintable(text, textIndex, part, partIndex) {
				part += PRINTABLE_CARET
			}

			currentLine.Text = append(
				currentLine.Text,
				core.NewFragment(part, f.Styles),
			)

			if partIndex < len(parts)-1 {
				lines = append(lines, currentLine)
				index++

				currentLine = core.FragmentLine(text.Padding)
				if c.index {
					currentLine.SetOrder(index)
				}
			}
		}
	}

	if len(currentLine.Text) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

func (c *TextArea) fixEmptyLines(lines []core.Line) []core.Line {
	for i, line := range lines {
		if line.Len() == 0 {
			styles := core.None
			if len(line.Text) > 0 {
				styles = line.Text[len(line.Text)-1].Styles
			}

			lines[i].Text = append(line.Text, core.NewFragment(EMPTY_LINE_FIX, styles))
		}
	}
	return lines
}

func (c *TextArea) isCaretPrintable(text core.Line, textIndex int, part string, partIndex int) bool {
	fragment := text.Text[textIndex]

	isCaret := len(part) == 0 && fragment.Styles.HasAny(core.Select)
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
