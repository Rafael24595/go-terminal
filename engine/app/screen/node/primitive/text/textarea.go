package text

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/line"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/textarea"
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/model/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/model/event"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
)

const NameArea = "text_area"

const ArgAreaBuffer param.Typed[[]rune] = "text_area_buffer"

var area_read_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit text"},
	},
	[]key.Action{
		key.ActionEnter,
	},
)

var area_write_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEsc:   {Code: []string{"ESC"}, Detail: "Save & Quit"},
		key.ActionEnter: {Code: []string{"RET"}, Detail: "New line"},
	},
	[]key.Action{
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
	buffer    *buffer.RuneBuffer
	clipboard *buffer.Clipboard
	caret     *input.TextCursor
}

func NewArea() *TextArea {
	runeBuffer := buffer.NewRuneBuffer().
		PushRules(rule.Full...)

	return &TextArea{
		reference: NameArea,
		history:   event.NewTextEventService(),
		writeMode: false,
		indexMode: false,
		buffer:    runeBuffer,
		clipboard: buffer.NewClipboard(),
		caret:     input.NewTextCursor(false),
	}
}

func (n *TextArea) SetName(name string) *TextArea {
	n.reference = name
	return n
}

func (n *TextArea) SetBuffer(buffer *buffer.RuneBuffer) *TextArea {
	if buffer != nil {
		n.buffer = buffer
	}
	return n
}

func (n *TextArea) WriteMode() *TextArea {
	n.writeMode = true
	return n
}

func (n *TextArea) ReadMode() *TextArea {
	n.writeMode = false
	return n
}

func (n *TextArea) EnableBlinking() *TextArea {
	n.caret.EnableBlinking()
	return n
}

func (n *TextArea) DisableBlinking() *TextArea {
	n.caret.DisableBlinking()
	return n
}

func (n *TextArea) AddText(text string) *TextArea {
	n.buffer.Append([]rune(text))
	n.caret.MoveCaretTo(n.buffer.Buffer(), n.buffer.Size())
	return n
}

func (n *TextArea) ShowIndex() *TextArea {
	n.indexMode = true
	return n
}

func (n *TextArea) HideIndex() *TextArea {
	n.indexMode = false
	return n
}

func (n *TextArea) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Definition(n.definition).
		Update(n.update).
		View(n.view).
		ToNode()
}

func (n *TextArea) definition() screen.Definition {
	if n.writeMode {
		return area_write_definition
	}
	return area_read_definition
}

func (n *TextArea) update(stt *state.UIState, evnt screen.Event) screen.Result {
	stt.Pager.ForceShow = true

	if !n.writeMode {
		return n.updateRead(stt, evnt)
	}
	return n.updateWrite(stt, evnt)
}

func (n *TextArea) updateRead(stt *state.UIState, evt screen.Event) screen.Result {
	ky := evt.Key

	switch ky.Code {
	case key.ActionEnter:
		n.writeMode = true
	}

	return screen.ResultFromUIState(stt)
}

func (n *TextArea) updateWrite(stt *state.UIState, evt screen.Event) screen.Result {
	ky := evt.Key

	switch ky.Code {
	case key.ActionEsc:
		n.writeMode = false
		return screen.ResultFromUIState(stt)

	case key.ActionHome:
		return n.moveHome(stt, evt)

	case key.ActionEnd:
		return n.moveEnd(stt, evt)

	case key.ActionArrowLeft:
		return n.moveBackward(stt, evt)

	case key.ActionArrowRight:
		return n.moveForward(stt, evt)

	case key.ActionEnter:
		ky = *key.NewKeyRune(ascii.ENTER_LF)

	case key.ActionArrowUp:
		return n.moveUp(stt, evt)

	case key.ActionArrowDown:
		return n.moveDown(stt, evt)
	}

	result := n.updateBuffer(stt, ky)

	state.PushParam(
		stt.Stack,
		n.reference,
		ArgAreaBuffer,
		n.buffer.Buffer(),
	)

	return result
}

func (n *TextArea) updateBuffer(state *state.UIState, ky key.Key) screen.Result {
	switch ky.Code {
	case key.ActionBackspace, key.ActionDeleteBackward:
		word := ky.Code == key.ActionDeleteBackward
		return n.deleteBackward(state, word)

	case key.ActionDelete, key.ActionDeleteForward:
		word := ky.Code == key.ActionDeleteForward
		return n.deleteForward(state, word)
	case key.CustomActionUndo, key.CustomActionRedo:
		return n.undoRedo(state, ky)

	case key.CustomActionCut, key.CustomActionCopy:
		cut := ky.Code == key.CustomActionCut
		return n.copyCut(state, cut)

	case key.CustomActionPaste:
		return n.paste(state)
	}

	return n.pushRune(state, ky)
}

func (n *TextArea) pushRune(state *state.UIState, ky key.Key) screen.Result {
	start, end, fixEnd := n.insertSelection()

	insert, delete := n.buffer.ReplaceWithRules([]rune{ky.Rune}, start, end)
	n.history.PushEvent(event.Insert, start, fixEnd, string(delete), string(insert))

	position := start + offset.Offset(len(insert))
	n.caret.MoveCaretTo(n.buffer.Buffer(), position)

	return screen.ResultFromUIState(state)
}

func (n *TextArea) undoRedo(state *state.UIState, ky key.Key) screen.Result {
	result := screen.ResultFromUIState(state)

	var delta *delta.Delta
	switch ky.Code {
	case key.CustomActionUndo:
		delta = n.history.Undo()
	case key.CustomActionRedo:
		delta = n.history.Redo()
	default:
		assert.Unreachable("unsupported key code '%d'", ky.Code)
		delta = n.history.Redo()
	}

	if delta == nil {
		return result
	}

	n.buffer.ApplyDelta(delta)

	position := delta.Start + delta.Measure()
	n.caret.MoveCaretTo(n.buffer.Buffer(), position)

	return result
}

func (n *TextArea) copyCut(state *state.UIState, cut bool) screen.Result {
	result := screen.ResultFromUIState(state)

	if n.buffer.Empty() {
		return result
	}

	start := n.caret.SelectStart().Sub(1)
	end := n.caret.SelectEnd()

	n.clipboard.Put(n.buffer.Range(start, end))

	if cut {
		n.history.PushEvent(event.Cut, start, end, string(n.clipboard.Buffer()), "")
		n.buffer.Delete(start, end)
		n.caret.MoveCaretTo(n.buffer.Buffer(), start)
	}

	return result
}

func (n *TextArea) paste(state *state.UIState) screen.Result {
	start, end, fixEnd := n.insertSelection()

	insert, delete := n.buffer.Replace(n.clipboard.Buffer(), start, end)
	n.history.PushEvent(event.Paste, start, fixEnd, string(delete), string(insert))

	position := start + offset.Offset(len(insert))
	n.caret.MoveCaretTo(n.buffer.Buffer(), position)

	return screen.ResultFromUIState(state)
}

func (n *TextArea) moveHome(state *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(state)

	buffer := n.buffer.Buffer()

	if event.Key.Mod.HasAny(key.ModCtrl) {
		n.caret.MoveCaretTo(buffer, 0)
		return result
	}

	caret := runes.BackwardIndexWithLimit(buffer, runes.NextLineRunes, n.caret.Caret())

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)

	return result
}

func (n *TextArea) moveEnd(state *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(state)

	buffer := n.buffer.Buffer()

	if event.Key.Mod.HasAny(key.ModCtrl) {
		n.caret.MoveCaretTo(buffer, n.buffer.Size())
		return result
	}

	caret := runes.ForwardIndexWithLimit(buffer, runes.NextLineRunes, n.caret.Caret())

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)

	return result
}

func (n *TextArea) moveUp(state *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(state)

	buffer := n.buffer.Buffer()

	start := n.caret.Caret()
	distance := line.DistanceFromLF(buffer, start)

	prevLineStart, ok := line.FindPrevLineStart(buffer, start)
	if !ok {
		if event.Key.Mod.HasAny(key.ModShift) {
			n.caret.MoveSelectTo(buffer, 0, n.caret.Anchor())
			return result
		}

		n.caret.MoveCaretTo(buffer, 0)
		return result
	}

	position := line.ClampToLine(buffer, prevLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		n.caret.MoveSelectTo(buffer, position, n.caret.Anchor())
	} else {
		n.caret.MoveCaretTo(buffer, position)
	}

	return result
}

func (n *TextArea) moveDown(state *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(state)

	buffer := n.buffer.Buffer()
	size := n.buffer.Size()

	start := n.caret.Caret()
	distance := line.DistanceFromLF(buffer, start)

	nextLineStart, ok := line.FindNextLineStart(buffer, start)
	if !ok {
		if event.Key.Mod.HasAny(key.ModShift) {
			n.caret.MoveSelectTo(buffer, size, n.caret.Anchor())
			return result
		}

		n.caret.MoveCaretTo(buffer, size)
		return result
	}

	position := line.ClampToLine(buffer, nextLineStart, distance)

	if event.Key.Mod.HasAny(key.ModShift) {
		n.caret.MoveSelectTo(buffer, position, n.caret.Anchor())
	} else {
		n.caret.MoveCaretTo(buffer, position)
	}

	return result
}

func (n *TextArea) moveBackward(state *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(state)

	buffer := n.buffer.Buffer()

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := n.caret.Caret().Sub(1)
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := n.caret.Caret().Sub(1)
		n.caret.MoveSelectTo(buffer, caret, anchor)
		return result
	}

	caret := runes.BackwardIndex(buffer, runes.NextWordRunes, n.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)
	return result
}

func (n *TextArea) moveForward(state *state.UIState, event screen.Event) screen.Result {
	result := screen.ResultFromUIState(state)

	buffer := n.buffer.Buffer()
	size := n.buffer.Size()

	if event.Key.Mod.HasNone(key.ModShift, key.ModCtrl) {
		caret := min(size, n.caret.Caret()+1)
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	anchor := n.caret.Anchor()
	if event.Key.Mod.HasNone(key.ModCtrl) {
		caret := min(size, n.caret.Caret()+1)
		n.caret.MoveSelectTo(buffer, caret, anchor)
		return result
	}

	caret := runes.ForwardIndex(buffer, runes.NextWordRunes, n.caret.Caret())
	if event.Key.Mod.HasNone(key.ModShift) {
		n.caret.MoveCaretTo(buffer, caret)
		return result
	}

	n.caret.MoveSelectTo(buffer, caret, anchor)
	return result
}

func (n *TextArea) deleteBackward(state *state.UIState, word bool) screen.Result {
	result := screen.ResultFromUIState(state)

	if n.buffer.Empty() {
		return result
	}

	start := n.caret.SelectStart()

	if word {
		start = runes.BackwardIndex(n.buffer.Buffer(), runes.NextWordRunes, start)
	} else {
		start = start.Sub(1)
	}

	end := n.caret.SelectEnd()

	delete := n.buffer.Delete(start, end)
	n.history.PushEvent(event.DeleteBackward, start, end, string(delete), "")

	n.caret.MoveCaretTo(n.buffer.Buffer(), start)
	return result
}

func (n *TextArea) deleteForward(state *state.UIState, word bool) screen.Result {
	result := screen.ResultFromUIState(state)

	if n.buffer.Empty() {
		return result
	}

	end := n.caret.SelectEnd()

	if word {
		end = runes.ForwardIndex(n.buffer.Buffer(), runes.NextWordRunes, end)
	} else {
		end = min(n.buffer.Size(), end+1)
	}

	start := n.caret.SelectStart().Sub(1)

	delete := n.buffer.Delete(start, end)
	n.history.PushEvent(event.DeleteForward, start, end, string(delete), "")

	n.caret.MoveCaretTo(n.buffer.Buffer(), start)
	return result
}

func (n *TextArea) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	predicate, textarea, needsPulse := n.viewSources()

	vm.Kernel.Push(
		textarea.ToUnit(),
	)

	vm.Pager.SetPredicate(predicate)
	vm.Behavior.NeedsPulse = needsPulse

	return *vm
}

func (n *TextArea) viewSources() (
	pager.Predicate,
	*textarea.TextAreaUnit,
	bool,
) {
	predicate := pager.PredicatePage()
	if n.writeMode {
		predicate = pager.PredicateFocus()
	}

	textarea := textarea.New(n.buffer.Facade(), n.caret).
		WriteMode(n.writeMode).
		IndexMode(n.indexMode)

	needsPulse := n.needsPulse()

	return predicate, textarea, needsPulse
}

func (n *TextArea) needsPulse() bool {
	return n.writeMode && n.caret.IsBlinking()
}

func (n *TextArea) insertSelection() (offset.Offset, offset.Offset, offset.Offset) {
	start := n.caret.SelectStart()
	end := n.caret.SelectEnd()

	if start != end {
		return start.Sub(1), end, end + 1
	}

	return start, end, end
}
