package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/core/event"
	"github.com/Rafael24595/go-terminal/engine/core/input"
	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper/line"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	text_helper "github.com/Rafael24595/go-terminal/engine/helper/text"

	drawable_line "github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/textarea"
)

const default_text_area_name = "TextArea"

var text_area_write_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionAll)...,
)

var text_area_read_definition = screen.DefinitionFromKeys(
	key.NewKeysCode(key.ActionEnter)...,
)

type TextArea struct {
	reference string
	history   *event.TextEventService
	writeMode bool
	indexMode bool
	title     []text.Line
	buffer    []rune
	caret     *input.Cursor
}

func NewTextArea() *TextArea {
	return &TextArea{
		reference: default_text_area_name,
		history:   event.NewTextEventService(),
		writeMode: false,
		title:     make([]text.Line, 0),
		caret:     input.NewCursor(false),
		buffer:    make([]rune, 0),
	}
}

func (c *TextArea) SetName(name string) *TextArea {
	c.reference = name
	return c
}

func (c *TextArea) WriteMode() *TextArea {
	c.writeMode = true
	return c
}

func (c *TextArea) ReadMode() *TextArea {
	c.writeMode = false
	return c
}

func (c *TextArea) EnableBlinking() *TextArea {
	c.caret.EnableBlinking()
	return c
}

func (c *TextArea) DisableBlinking() *TextArea {
	c.caret.DisableBlinking()
	return c
}

func (c *TextArea) AddTitle(title ...text.Line) *TextArea {
	c.title = append(c.title, title...)
	return c
}

func (c *TextArea) AddText(text string) *TextArea {
	c.buffer = append(c.buffer, []rune(text)...)
	c.caret.MoveCaretTo(c.buffer, uint(len(c.buffer)))
	return c
}

func (c *TextArea) ShowIndex() *TextArea {
	c.indexMode = true
	return c
}

func (c *TextArea) HideIndex() *TextArea {
	c.indexMode = false
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
	if c.writeMode {
		return text_area_write_definition
	}
	return text_area_read_definition
}

func (c *TextArea) update(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	state.Pager.ShowPage = true

	if !c.writeMode {
		return c.updateRead(state, evnt)
	}
	return c.updateWrite(state, evnt)
}

func (c *TextArea) updateRead(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEnter:
		c.writeMode = true
	}

	return screen.ScreenResultFromUIState(state)
}

func (c *TextArea) updateWrite(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	ky := evnt.Key

	switch ky.Code {
	case key.ActionEsc:
		c.writeMode = false
		return screen.ScreenResultFromUIState(state)
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

	return c.pushRune(state, ky)
}

func (c *TextArea) pushRune(state *state.UIState, ky key.Key) screen.ScreenResult {
	end := c.caret.SelectEnd()

	start := c.caret.SelectStart()
	fixEnd := end
	if start != end {
		start = math.SubClampZero(start, 1)
		fixEnd += 1
	}

	text := text_helper.FullTextTransformer.Apply(ky, start, end, c.buffer)

	c.history.PushEvent(event.Insert, start, fixEnd, string(c.buffer[start:end]), string(text))

	c.buffer = runes.AppendRange(c.buffer, text, start, end)

	position := start + uint(len(text))
	c.caret.MoveCaretTo(c.buffer, position)

	return screen.ScreenResultFromUIState(state)
}

func (c *TextArea) undoRedo(state *state.UIState, ky key.Key) screen.ScreenResult {
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
	c.caret.MoveCaretTo(c.buffer, delta.Start+uint(len(newTextRunes)))

	return result
}

func (c *TextArea) moveHome(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasAny(key.ModCtrl) {
		c.caret.MoveCaretTo(c.buffer, 0)
		return result
	}

	caret := runes.BackwardIndexWithLimit(c.buffer, runes.NextLineRunes, c.caret.Caret())

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(c.buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(c.buffer, caret, anchor)

	return result
}

func (c *TextArea) moveEnd(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasAny(key.ModCtrl) {
		c.caret.MoveCaretTo(c.buffer, uint(len(c.buffer)))
		return result
	}

	caret := runes.ForwardIndexWithLimit(c.buffer, runes.NextLineRunes, c.caret.Caret())

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(c.buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(c.buffer, caret, anchor)

	return result
}

func (c *TextArea) moveUp(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	start := c.caret.Caret()
	distance := line.DistanceFromLF(c.buffer, int(start))

	prevLineStart := line.FindPrevLineStart(c.buffer, int(start))
	if prevLineStart == -1 {
		if event.Key.Mod.HasAny(key.ModShift) {
			c.caret.MoveSelectTo(c.buffer, 0, c.caret.Anchor())
			return result
		}

		c.caret.MoveCaretTo(c.buffer, 0)
		return result
	}

	position := line.ClampToLine(c.buffer, prevLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		c.caret.MoveSelectTo(c.buffer, uint(position), c.caret.Anchor())
	} else {
		c.caret.MoveCaretTo(c.buffer, uint(position))
	}

	return result
}

func (c *TextArea) moveDown(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	start := c.caret.Caret()
	distance := line.DistanceFromLF(c.buffer, int(start))

	nextLineStart := line.FindNextLineStart(c.buffer, int(start))
	if nextLineStart == -1 {
		if event.Key.Mod.HasAny(key.ModShift) {
			c.caret.MoveSelectTo(c.buffer, uint(len(c.buffer)), c.caret.Anchor())
			return result
		}

		c.caret.MoveCaretTo(c.buffer, uint(len(c.buffer)))
		return result
	}

	position := line.ClampToLine(c.buffer, nextLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		c.caret.MoveSelectTo(c.buffer, uint(position), c.caret.Anchor())
	} else {
		c.caret.MoveCaretTo(c.buffer, uint(position))
	}

	return result
}

func (c *TextArea) moveBackward(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := math.SubClampZero(c.caret.Caret(), 1)
		c.caret.MoveCaretTo(c.buffer, caret)
		return result
	}

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := math.SubClampZero(c.caret.Caret(), 1)
		c.caret.MoveSelectTo(c.buffer, caret, anchor)
		return result
	}

	caret := runes.BackwardIndex(c.buffer, runes.NextWordRunes, c.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(c.buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(c.buffer, caret, anchor)
	return result
}

func (c *TextArea) moveForward(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := min(uint(len(c.buffer)), c.caret.Caret()+1)
		c.caret.MoveCaretTo(c.buffer, caret)
		return result
	}

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := min(uint(len(c.buffer)), c.caret.Caret()+1)
		c.caret.MoveSelectTo(c.buffer, caret, anchor)
		return result
	}

	caret := runes.ForwardIndex(c.buffer, runes.NextWordRunes, c.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(c.buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(c.buffer, caret, anchor)
	return result
}

func (c *TextArea) deleteBackward(state *state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if len(c.buffer) == 0 {
		return result
	}

	start := c.caret.SelectStart()

	if word {
		start = runes.BackwardIndex(c.buffer, runes.NextWordRunes, start)
	} else {
		start = math.SubClampZero(start, 1)
	}

	end := c.caret.SelectEnd()

	c.history.PushEvent(event.DeleteBackward, start, end, string(c.buffer[start:end]), "")

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)

	c.caret.MoveCaretTo(c.buffer, start)
	return result
}

func (c *TextArea) deleteForward(state *state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if len(c.buffer) == 0 {
		return result
	}

	end := c.caret.SelectEnd()

	if word {
		end = runes.ForwardIndex(c.buffer, runes.NextWordRunes, end)
	} else {
		end = min(uint(len(c.buffer)), end+1)
	}

	start := math.SubClampZero(c.caret.SelectStart(), 1)

	c.history.PushEvent(event.DeleteForward, start, end, string(c.buffer[start:end]), "")

	c.buffer = append(c.buffer[:start], c.buffer[end:]...)

	c.caret.MoveCaretTo(c.buffer, start)
	return result
}

func (c *TextArea) view(stt state.UIState) core.ViewModel {
	strategy := state.NewPagePager()
	if c.writeMode {
		strategy = state.NewFocusPager()
	}

	textarea := textarea.NewTextAreaDrawable(c.buffer, *c.caret).
		WriteMode(c.writeMode).
		IndexMode(c.indexMode)

	vm := core.ViewModelFromUIState(stt)

	vm.Header.Shift(
		drawable_line.EagerDrawableFromLines(c.title...),
	)
	vm.Lines.Shift(
		textarea.ToDrawable(),
	)

	vm.SetStrategy(strategy)

	return *vm
}
