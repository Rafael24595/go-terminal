package primitive

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/app/pager"
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/helper/line"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/block"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/textarea"
	"github.com/Rafael24595/go-terminal/engine/model/buffer"
	"github.com/Rafael24595/go-terminal/engine/model/delta"
	"github.com/Rafael24595/go-terminal/engine/model/event"
	"github.com/Rafael24595/go-terminal/engine/model/help"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/model/param"
	"github.com/Rafael24595/go-terminal/engine/render/text"

	text_transformer "github.com/Rafael24595/go-terminal/engine/helper/text"
)

const default_text_area_name = "TextArea"

const ArgTextAreaBuffer param.Typed[[]rune] = "text_area_buffer"

var text_area_read_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionEsc:   {Code: []string{"ESC"}, Detail: "Exit/Back"},
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit text"},
	},
	[]key.KeyAction{
		key.ActionEnter,
	},
)

var text_area_write_definition = screen.NewDefinitionSources(
	map[key.KeyAction]help.HelpField{
		key.ActionEsc:   {Code: []string{"ESC"}, Detail: "Save & Quit"},
		key.ActionEnter: {Code: []string{"RET"}, Detail: "New line"},
	},
	[]key.KeyAction{
		key.ActionEsc,
		key.ActionHome,
		key.ActionEnd,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionBackspace,
		key.ActionDeleteBackward,
		key.ActionDelete,
		key.ActionDeleteForward,
		key.ActionEnter,
		key.ActionArrowUp,
		key.ActionArrowDown,
		key.CustomActionUndo,
		key.CustomActionRedo,
		key.CustomActionCut,
		key.CustomActionCopy,
		key.CustomActionPaste,
		key.ActionRune,
	},
)

type TextArea struct {
	reference string
	history   *event.TextEventService
	writeMode bool
	indexMode bool
	title     []text.Line
	buffer    *buffer.RuneBuffer
	clipboard *buffer.Clipboard
	caret     *input.TextCursor
}

func NewTextArea() *TextArea {
	runeBuffer := buffer.NewRuneBuffer().
		Transformer(text_transformer.FullTextTransformer)

	return &TextArea{
		reference: default_text_area_name,
		history:   event.NewTextEventService(),
		writeMode: false,
		indexMode: false,
		title:     make([]text.Line, 0),
		buffer:    runeBuffer,
		clipboard: buffer.NewClipboard(),
		caret:     input.NewTextCursor(false),
	}
}

func (c *TextArea) SetName(name string) *TextArea {
	c.reference = name
	return c
}

func (c *TextArea) SetBuffer(buffer *buffer.RuneBuffer) *TextArea {
	if buffer != nil {
		c.buffer = buffer
	}
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
	c.buffer.Append([]rune(text))
	c.caret.MoveCaretTo(c.buffer.Buffer(), c.buffer.Size())
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
	screen := screen.Screen{
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}

	return screen.SetName(c.reference).
		StackFromName()
}

func (c *TextArea) definitionSource() screen.DefinitionSources {
	if c.writeMode {
		return text_area_write_definition
	}
	return text_area_read_definition
}

func (c *TextArea) definition() screen.Definition {
	return c.definitionSource().Definition
}

func (c *TextArea) update(state *state.UIState, evnt screen.ScreenEvent) screen.ScreenResult {
	state.Pager.ForceShow = true

	if !c.writeMode {
		return c.updateRead(state, evnt)
	}
	return c.updateWrite(state, evnt)
}

func (c *TextArea) updateRead(stt *state.UIState, evt screen.ScreenEvent) screen.ScreenResult {
	ky := evt.Key

	switch ky.Code {
	case key.ActionEnter:
		c.writeMode = true
	}

	return screen.ScreenResultFromUIState(stt)
}

func (c *TextArea) updateWrite(stt *state.UIState, evt screen.ScreenEvent) screen.ScreenResult {
	ky := evt.Key

	switch ky.Code {
	case key.ActionEsc:
		c.writeMode = false
		return screen.ScreenResultFromUIState(stt)

	case key.ActionHome:
		return c.moveHome(stt, evt)

	case key.ActionEnd:
		return c.moveEnd(stt, evt)

	case key.ActionArrowLeft:
		return c.moveBackward(stt, evt)

	case key.ActionArrowRight:
		return c.moveForward(stt, evt)

	case key.ActionEnter:
		ky = *key.NewKeyRune(key.ENTER_LF)

	case key.ActionArrowUp:
		return c.moveUp(stt, evt)

	case key.ActionArrowDown:
		return c.moveDown(stt, evt)
	}

	result := c.updateBuffer(stt, ky)

	state.PushParam(
		stt.Stack,
		c.reference,
		ArgTextAreaBuffer,
		c.buffer.Buffer(),
	)

	return result
}

func (c *TextArea) updateBuffer(state *state.UIState, ky key.Key) screen.ScreenResult {
	switch ky.Code {
	case key.ActionBackspace, key.ActionDeleteBackward:
		word := ky.Code == key.ActionDeleteBackward
		return c.deleteBackward(state, word)

	case key.ActionDelete, key.ActionDeleteForward:
		word := ky.Code == key.ActionDeleteForward
		return c.deleteForward(state, word)
	case key.CustomActionUndo, key.CustomActionRedo:
		return c.undoRedo(state, ky)

	case key.CustomActionCut, key.CustomActionCopy:
		cut := ky.Code == key.CustomActionCut
		return c.copyCut(state, cut)

	case key.CustomActionPaste:
		return c.paste(state)
	}

	return c.pushRune(state, ky)
}

func (c *TextArea) pushRune(state *state.UIState, ky key.Key) screen.ScreenResult {
	start, end, fixEnd := c.insertSelection()

	insert, delete := c.buffer.TransformAndReplace([]rune{ky.Rune}, start, end)
	c.history.PushEvent(event.Insert, start, fixEnd, string(delete), string(insert))

	position := start + uint(len(insert))
	c.caret.MoveCaretTo(c.buffer.Buffer(), position)

	return screen.ScreenResultFromUIState(state)
}

func (c *TextArea) undoRedo(state *state.UIState, ky key.Key) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	var delta *delta.Delta
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

	c.buffer.ApplyDelta(delta)

	position := delta.Start + delta.Measure()
	c.caret.MoveCaretTo(c.buffer.Buffer(), position)

	return result
}

func (c *TextArea) copyCut(state *state.UIState, cut bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if c.buffer.Empty() {
		return result
	}

	start := math.SubClampZero(c.caret.SelectStart(), 1)
	end := c.caret.SelectEnd()

	c.clipboard.Put(c.buffer.Range(start, end))

	if cut {
		c.history.PushEvent(event.Cut, start, end, string(c.clipboard.Buffer()), "")
		c.buffer.Delete(start, end)
		c.caret.MoveCaretTo(c.buffer.Buffer(), start)
	}

	return result
}

func (c *TextArea) paste(state *state.UIState) screen.ScreenResult {
	start, end, fixEnd := c.insertSelection()

	insert, delete := c.buffer.Replace(c.clipboard.Buffer(), start, end)
	c.history.PushEvent(event.Paste, start, fixEnd, string(delete), string(insert))

	position := start + uint(len(insert))
	c.caret.MoveCaretTo(c.buffer.Buffer(), position)

	return screen.ScreenResultFromUIState(state)
}

func (c *TextArea) moveHome(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	buffer := c.buffer.Buffer()

	if event.Key.Mod.HasAny(key.ModCtrl) {
		c.caret.MoveCaretTo(buffer, 0)
		return result
	}

	caret := runes.BackwardIndexWithLimit(buffer, runes.NextLineRunes, c.caret.Caret())

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(buffer, caret, anchor)

	return result
}

func (c *TextArea) moveEnd(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	buffer := c.buffer.Buffer()

	if event.Key.Mod.HasAny(key.ModCtrl) {
		c.caret.MoveCaretTo(buffer, c.buffer.Size())
		return result
	}

	caret := runes.ForwardIndexWithLimit(buffer, runes.NextLineRunes, c.caret.Caret())

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(buffer, caret, anchor)

	return result
}

func (c *TextArea) moveUp(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	buffer := c.buffer.Buffer()

	start := c.caret.Caret()
	distance := line.DistanceFromLF(buffer, int(start))

	prevLineStart := line.FindPrevLineStart(buffer, int(start))
	if prevLineStart == -1 {
		if event.Key.Mod.HasAny(key.ModShift) {
			c.caret.MoveSelectTo(buffer, 0, c.caret.Anchor())
			return result
		}

		c.caret.MoveCaretTo(buffer, 0)
		return result
	}

	position := line.ClampToLine(buffer, prevLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		c.caret.MoveSelectTo(buffer, uint(position), c.caret.Anchor())
	} else {
		c.caret.MoveCaretTo(buffer, uint(position))
	}

	return result
}

func (c *TextArea) moveDown(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	buffer := c.buffer.Buffer()
	size := c.buffer.Size()

	start := c.caret.Caret()
	distance := line.DistanceFromLF(buffer, int(start))

	nextLineStart := line.FindNextLineStart(buffer, int(start))
	if nextLineStart == -1 {
		if event.Key.Mod.HasAny(key.ModShift) {
			c.caret.MoveSelectTo(buffer, size, c.caret.Anchor())
			return result
		}

		c.caret.MoveCaretTo(buffer, size)
		return result
	}

	position := line.ClampToLine(buffer, nextLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		c.caret.MoveSelectTo(buffer, uint(position), c.caret.Anchor())
	} else {
		c.caret.MoveCaretTo(buffer, uint(position))
	}

	return result
}

func (c *TextArea) moveBackward(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	buffer := c.buffer.Buffer()

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := math.SubClampZero(c.caret.Caret(), 1)
		c.caret.MoveCaretTo(buffer, caret)
		return result
	}

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := math.SubClampZero(c.caret.Caret(), 1)
		c.caret.MoveSelectTo(buffer, caret, anchor)
		return result
	}

	caret := runes.BackwardIndex(buffer, runes.NextWordRunes, c.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(buffer, caret, anchor)
	return result
}

func (c *TextArea) moveForward(state *state.UIState, event screen.ScreenEvent) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	buffer := c.buffer.Buffer()
	size := c.buffer.Size()

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := min(size, c.caret.Caret()+1)
		c.caret.MoveCaretTo(buffer, caret)
		return result
	}

	anchor := c.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := min(size, c.caret.Caret()+1)
		c.caret.MoveSelectTo(buffer, caret, anchor)
		return result
	}

	caret := runes.ForwardIndex(buffer, runes.NextWordRunes, c.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		c.caret.MoveCaretTo(buffer, caret)
		return result
	}

	c.caret.MoveSelectTo(buffer, caret, anchor)
	return result
}

func (c *TextArea) deleteBackward(state *state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if c.buffer.Empty() {
		return result
	}

	start := c.caret.SelectStart()

	if word {
		start = runes.BackwardIndex(c.buffer.Buffer(), runes.NextWordRunes, start)
	} else {
		start = math.SubClampZero(start, 1)
	}

	end := c.caret.SelectEnd()

	delete := c.buffer.Delete(start, end)
	c.history.PushEvent(event.DeleteBackward, start, end, string(delete), "")

	c.caret.MoveCaretTo(c.buffer.Buffer(), start)
	return result
}

func (c *TextArea) deleteForward(state *state.UIState, word bool) screen.ScreenResult {
	result := screen.ScreenResultFromUIState(state)

	if c.buffer.Empty() {
		return result
	}

	end := c.caret.SelectEnd()

	if word {
		end = runes.ForwardIndex(c.buffer.Buffer(), runes.NextWordRunes, end)
	} else {
		end = min(c.buffer.Size(), end+1)
	}

	start := math.SubClampZero(c.caret.SelectStart(), 1)

	delete := c.buffer.Delete(start, end)
	c.history.PushEvent(event.DeleteForward, start, end, string(delete), "")

	c.caret.MoveCaretTo(c.buffer.Buffer(), start)
	return result
}

func (c *TextArea) view(stt state.UIState) viewmodel.ViewModel {
	source := c.definitionSource()

	predicate := pager.PredicatePage()
	if c.writeMode {
		predicate = pager.PredicateFocus()
	}

	textarea := textarea.NewTextAreaDrawable(c.buffer.Facade(), c.caret).
		WriteMode(c.writeMode).
		IndexMode(c.indexMode)

	vm := viewmodel.ViewModelFromUIState(stt)

	vm.Header.Push(
		block.BlockDrawableFromLines(c.title...),
	)

	code := c.mainDrawableCode()
	vm.Kernel.Push(
		textarea.ToDrawable().SetCode(code),
	)

	vm.Pager.SetPredicate(predicate)

	vm.Helper.Shift(
		key.ActionsToHelpWithOverride(
			source.Overrides, source.Actions...,
		)...,
	)

	return *vm
}

func (c *TextArea) mainDrawableCode() string {
	return "main_" + c.reference
}

func (c *TextArea) insertSelection() (uint, uint, uint) {
	start := c.caret.SelectStart()
	end := c.caret.SelectEnd()

	if start != end {
		return math.SubClampZero(start, 1), end, end + 1
	}

	return start, end, end
}
