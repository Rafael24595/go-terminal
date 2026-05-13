package selection

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Result struct {
	Frags []text.Fragment
	End   offset.Offset
}

type Renderer struct {
	buffer []rune
	start  offset.Offset
	end    offset.Offset
	blink  style.Atom
}

func NewRenderer(
	buffer []rune,
	start, end offset.Offset,
	blink ...style.Atom,
) Renderer {
	return Renderer{
		buffer: buffer,
		start:  start,
		end:    end,
		blink:  style.MergeAtom(blink...),
	}
}

func (r Renderer) selection() []rune {
	return r.buffer[r.start:r.end]
}

func (r Renderer) Resolve(caret *input.TextCursor) Result {
	selection := r.selection()
	if len(selection) == 0 {
		return r.resolveEmpty()
	}

	if caret.Caret() != caret.Anchor() && r.end == caret.Anchor() {
		return r.resolveBackward()
	}

	return r.resolveForward()
}

func (r Renderer) resolveBackward() Result {
	frags := make([]text.Fragment, 0, 2)
	focusAtom := style.AtmFocus

	selection := r.selection()
	if r.start > 0 && selection[0] == ascii.ENTER_LF {
		focusAtom = style.AtmNone

		frags = append(frags,
			frag(marker.PrintableCaretRunes, r.blink, style.AtmFocus),
		)
	}

	frags = append(frags,
		frag(selection, r.blink, focusAtom),
	)

	return Result{
		Frags: frags,
		End:   r.end,
	}
}

func (r Renderer) resolveForward() Result {
	selection := r.selection()
	if selection[len(selection)-1] == ascii.ENTER_LF {
		return r.resolveForwardEnter()
	}

	return r.resolveForwardNonEnter()
}

func (r Renderer) resolveForwardNonEnter() Result {
	frags := make([]text.Fragment, 0, 3)

	selection := r.selection()
	if len(selection) > 1 {
		frags = append(frags,
			frag(selection[:len(selection)-1], r.blink),
		)
	}

	frags = append(frags,
		frag(selection[len(selection)-1:], r.blink, style.AtmFocus),
	)

	return Result{
		Frags: frags,
		End:   r.end,
	}
}

func (r Renderer) resolveForwardEnter() Result {
	frags := make([]text.Fragment, 0, 3)

	selection := r.selection()
	if len(selection) == 1 {
		frags = append(frags,
			frag(marker.PrintableCaretRunes, r.blink),
		)
	}

	footer, nextEnd := r.resolveEnterFooter()

	frags = append(frags,
		frag(selection, r.blink),
		frag(footer, r.blink, style.AtmFocus),
	)

	return Result{
		Frags: frags,
		End:   nextEnd,
	}
}

func (r Renderer) resolveEnterFooter() ([]rune, offset.Offset) {
	if int(r.end) >= len(r.buffer) {
		return marker.PrintableCaretRunes, r.end
	}

	if r.buffer[r.end] == ascii.ENTER_LF {
		return marker.PrintableCaretRunes, r.end
	}

	return r.buffer[r.end : r.end+1], r.end + 1
}

func (r Renderer) resolveEmpty() Result {
	assert.Unreachable("selection should have at least one character")

	frags := []text.Fragment{
		*text.EmptyFragment().AddAtom(style.AtmFocus),
	}

	return Result{
		Frags: frags,
		End:   r.end,
	}
}

func frag(runes []rune, atoms ...style.Atom) text.Fragment {
	return *text.FragmentFromRunes(runes).
		AddAtom(atoms...)
}
