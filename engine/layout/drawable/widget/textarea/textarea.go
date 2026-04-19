package textarea

import (
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameTextAreaDrawable = "TextAreaDrawable"

type TextAreaDrawable struct {
	loaded     bool
	lazyLoaded bool
	writeMode  bool
	indexMode  bool
	buffer     []rune
	caret      *input.TextCursor
	drawable   drawable.Drawable
}

func NewTextAreaDrawable(buffer []rune, caret *input.TextCursor) *TextAreaDrawable {
	clone := make([]rune, len(buffer))
	copy(clone, buffer)

	return &TextAreaDrawable{
		loaded:     false,
		lazyLoaded: false,
		writeMode:  false,
		indexMode:  false,
		buffer:     clone,
		caret:      caret,
		drawable:   drawable.Drawable{},
	}
}

func TextAreaDrawableFromData(buffer []rune, caret *input.TextCursor) drawable.Drawable {
	return NewTextAreaDrawable(buffer, caret).ToDrawable()
}

func (d *TextAreaDrawable) WriteMode(writeMode bool) *TextAreaDrawable {
	d.writeMode = writeMode
	return d
}

func (d *TextAreaDrawable) IndexMode(indexMode bool) *TextAreaDrawable {
	d.indexMode = indexMode
	return d
}

func (d *TextAreaDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: NameTextAreaDrawable,
		Code: d.drawable.Code,
		Tags: d.drawable.Tags,
		Init: d.init,
		Wipe: d.wipe,
		Draw: d.draw,
	}
}

func (d *TextAreaDrawable) init() {
	d.loaded = true
	d.lazyLoaded = false
}

func (d *TextAreaDrawable) lazyInit(size terminal.Winsize) {
	if d.lazyLoaded {
		return
	}

	d.lazyLoaded = true

	start := math.SubClampZero(d.caret.SelectStart(), 1)
	end := d.caret.SelectEnd()

	if len(d.buffer) == 0 {
		d.buffer = append(d.buffer, marker.PrintableCaretRunes...)
		start = 0
		end = 1
	}

	txt := text.EmptyLine().
		SetSpec(style.SpecFromKind(style.SpcKindPaddingRight))

	fragments := d.resolveFragments(d.buffer, start, end)
	txt.Text = append(txt.Text, fragments...)

	lines := d.normalizeLinesEnd(*txt)
	lines = d.fixEmptyLines(size, lines)

	drawable := line.LineDrawableFromLines(lines...)
	drawable.Init()

	d.drawable = drawable
}

func (d *TextAreaDrawable) wipe() {
	d.lazyLoaded = false

	if d.drawable.Wipe == nil {
		return
	}

	d.drawable.Wipe()
}

func (d *TextAreaDrawable) resolveFragments(renderBuffer []rune, start uint, end uint) []text.Fragment {
	frags := make([]text.Fragment, 0, 6)

	if int(start) > 0 {
		frags = append(frags, *text.NewFragment(string(renderBuffer[:start])))
	}

	var selection []text.Fragment
	if d.caret.Caret() != d.caret.Anchor() && end == d.caret.Anchor() {
		selection, start, end = d.resolveBackwardSelection(renderBuffer, start, end)
	} else {
		selection, start, end = d.resolveForwardSelection(renderBuffer, start, end)
	}

	frags = append(frags, selection...)

	if int(end) < len(renderBuffer) {
		frags = append(frags, *text.NewFragment(string(renderBuffer[end:])))
	}

	return frags
}

func (d *TextAreaDrawable) resolveBackwardSelection(renderBuffer []rune, start uint, end uint) ([]text.Fragment, uint, uint) {
	selection := renderBuffer[start:end]
	caretAtom := d.BlinkStyle()

	frags := make([]text.Fragment, 0, 2)

	focusAtom := style.AtmFocus

	if len(selection) == 0 {
		assert.Unreachable("selection should have at least one character")

		selectFrag := text.EmptyFragment().
			AddAtom(style.AtmFocus)

		frags = append(frags, *selectFrag)
		return frags, start, end
	}

	if int(start) > 0 && selection[0] == key.ENTER_LF {
		focusAtom = style.AtmNone

		headerFrag := text.FragmentFromRunes(marker.PrintableCaretRunes).
			AddAtom(caretAtom, style.AtmFocus)

		frags = append(frags, *headerFrag)
	}

	selectFrag := text.FragmentFromRunes(selection).
		AddAtom(caretAtom, focusAtom)

	frags = append(frags, *selectFrag)
	return frags, start, end
}

func (d *TextAreaDrawable) resolveForwardSelection(renderBuffer []rune, start uint, end uint) ([]text.Fragment, uint, uint) {
	selection := renderBuffer[start:end]
	caretAtom := d.BlinkStyle()

	frags := make([]text.Fragment, 0, 3)

	selectionSize := len(selection)
	if selectionSize == 0 {
		assert.Unreachable("selection should have at least one character")

		selectFrag := text.EmptyFragment().
			AddAtom(style.AtmFocus)

		frags = append(frags, *selectFrag)
		return frags, start, end
	}

	if selection[selectionSize-1] != key.ENTER_LF {
		headerFrag := text.FragmentFromRunes(selection[:selectionSize-1]).
			AddAtom(caretAtom)
		footerFrag := text.FragmentFromRunes(selection[selectionSize-1:]).
			AddAtom(caretAtom, style.AtmFocus)

		frags = append(frags, *headerFrag, *footerFrag)
		return frags, start, end
	}

	footer := marker.PrintableCaretRunes
	if int(end) < len(renderBuffer)-1 && renderBuffer[end+1] != key.ENTER_LF {
		footer = renderBuffer[end : end+1]
		end += 1
	}

	headerFrag := text.FragmentFromRunes(marker.PrintableCaretRunes).
		AddAtom(caretAtom)
	selectFrag := text.FragmentFromRunes(selection).
		AddAtom(caretAtom)
	footerFrag := text.FragmentFromRunes(footer).
		AddAtom(caretAtom, style.AtmFocus)

	frags = append(frags, *headerFrag, *selectFrag, *footerFrag)

	return frags, start, end
}

func (c *TextAreaDrawable) BlinkStyle() style.Atom {
	if !c.writeMode {
		return style.AtmNone
	}

	return c.caret.BlinkStyle()
}

func (d *TextAreaDrawable) normalizeLinesEnd(txt text.Line) []text.Line {
	lines := make([]text.Line, 0)

	index := uint16(1)

	currentLine := text.EmptyLine().SetSpec(txt.Spec)
	if d.indexMode {
		currentLine.SetOrder(index)
	}

	for _, f := range txt.Text {
		normalized := runes.NormalizeLineEnd(f.Text)

		parts := strings.Split(normalized, "\n")
		if len(parts) == 1 {
			currentLine.Text = append(
				currentLine.Text,
				*text.NewFragment(parts[0]).CopyMeta(&f),
			)

			continue
		}

		for partIndex, part := range parts {
			currentLine.Text = append(
				currentLine.Text,
				*text.NewFragment(part).CopyMeta(&f),
			)

			if partIndex >= len(parts)-1 {
				continue
			}

			lines = append(lines, *currentLine)
			index++

			currentLine = text.EmptyLine().SetSpec(txt.Spec)
			if d.indexMode {
				currentLine.SetOrder(index)
			}
		}
	}

	if len(currentLine.Text) > 0 {
		lines = append(lines, *currentLine)
	}

	return lines
}

func (d *TextAreaDrawable) fixEmptyLines(size terminal.Winsize, lines []text.Line) []text.Line {
	for i, line := range lines {
		if text.FragmentMeasure(int(size.Cols), line.Text...) != 0 {
			continue
		}

		styles := style.AtmNone
		if len(line.Text) > 0 {
			styles = line.Text[len(line.Text)-1].Atom
		}

		lines[i].Text = append(line.Text,
			*text.NewFragment(marker.DefaultPaddingText).
				AddAtom(styles),
		)
	}
	return lines
}
func (d *TextAreaDrawable) draw(size terminal.Winsize) ([]text.Line, bool) {
	assert.True(d.loaded, drawable.MessageInitialized)

	d.lazyInit(size)

	return d.drawable.Draw(size)
}
