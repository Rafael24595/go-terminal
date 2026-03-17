package textarea

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/model/input"
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/platform/assert"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const NameTextAreaDrawable = "TextAreaDrawable"

type TextAreaDrawable struct {
	initialized bool
	writeMode   bool
	indexMode   bool
	buffer      []rune
	caret       input.TextCursor
	drawable    drawable.Drawable
}

func NewTextAreaDrawable(buffer []rune, caret input.TextCursor) *TextAreaDrawable {
	clone := make([]rune, len(buffer))
	copy(clone, buffer)

	return &TextAreaDrawable{
		initialized: false,
		writeMode:   false,
		indexMode:   false,
		buffer:      clone,
		caret:       caret,
		drawable:    drawable.Drawable{},
	}
}

func TextAreaDrawableFromData(buffer []rune, caret input.TextCursor) drawable.Drawable {
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
		Init: d.init,
		Draw: d.draw,
	}
}

func (d *TextAreaDrawable) init(size terminal.Winsize) {
	d.initialized = true

	start := math.SubClampZero(d.caret.SelectStart(), 1)
	end := d.caret.SelectEnd()

	if len(d.buffer) == 0 {
		d.buffer = append(d.buffer, []rune(marker.PrintableCaretText)...)
		start = 0
		end = 1
	}

	txt := text.FragmentLine(style.SpecFromKind(style.SpcKindPaddingRight))

	beforeSelect := string(d.buffer[0:start])
	txt.Text = append(txt.Text, text.NewFragment(beforeSelect))

	onSelect := d.makeSelectedFragments(d.buffer, start, end)
	txt.Text = append(txt.Text, onSelect...)

	afterSelect := string(d.buffer[end:])
	if len(afterSelect) > 0 {
		txt.Text = append(txt.Text, text.NewFragment(afterSelect))
	}

	lines := d.normalizeLinesEnd(txt)
	lines = d.fixEmptyLines(lines)

	drawable := line.LazyDrawableFromLines(lines...)
	drawable.Init(size)

	d.drawable = drawable
}

func (d *TextAreaDrawable) makeSelectedFragments(renderBuffer []rune, start uint, end uint) []text.Fragment {
	onSelect := renderBuffer[start:end]

	selectAtom := style.AtmNone
	if d.writeMode {
		selectAtom = d.caret.BlinkStyle()
	}

	if d.caret.Caret() == d.caret.Anchor() {
		return []text.Fragment{
			text.NewFragment(string(onSelect)).
				AddAtom(selectAtom, style.AtmFocus),
		}
	}

	if end == d.caret.Anchor() {
		return []text.Fragment{
			text.NewFragment(string(onSelect[:1])).
				AddAtom(selectAtom, style.AtmFocus),
			text.NewFragment(string(onSelect[1:])).
				AddAtom(selectAtom),
		}
	}

	return []text.Fragment{
		text.NewFragment(string(onSelect[:len(onSelect)-1])).
			AddAtom(selectAtom),
		text.NewFragment(string(onSelect[len(onSelect)-1])).
			AddAtom(selectAtom, style.AtmFocus),
	}
}

func (d *TextAreaDrawable) normalizeLinesEnd(txt text.Line) []text.Line {
	lines := make([]text.Line, 0)

	index := uint16(1)

	currentLine := text.FragmentLine(txt.Spec)
	if d.indexMode {
		currentLine.SetOrder(index)
	}

	for textIndex, f := range txt.Text {
		normalized := runes.NormalizeLineEnd(f.Text)

		parts := strings.Split(normalized, "\n")
		if len(parts) == 1 {
			currentLine.Text = append(
				currentLine.Text,
				text.FragmentFrom(parts[0], f),
			)

			continue
		}

		for partIndex, part := range parts {
			if d.isCaretPrintable(txt, textIndex, part, partIndex) {
				part += marker.PrintableCaretText
			}

			currentLine.Text = append(
				currentLine.Text,
				text.FragmentFrom(part, f),
			)

			if partIndex >= len(parts)-1 {
				continue
			}

			lines = append(lines, currentLine)
			index++

			currentLine = text.FragmentLine(txt.Spec)
			if d.indexMode {
				currentLine.SetOrder(index)
			}
		}
	}

	if len(currentLine.Text) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

func (d *TextAreaDrawable) fixEmptyLines(lines []text.Line) []text.Line {
	for i, line := range lines {
		if text.LineFragmentsMeasure(line) != 0 {
			continue
		}

		styles := style.AtmNone
		if len(line.Text) > 0 {
			styles = line.Text[len(line.Text)-1].Atom
		}

		lines[i].Text = append(line.Text,
			text.NewFragment(marker.DefaultPaddingText).
				AddAtom(styles),
		)
	}
	return lines
}

func (d *TextAreaDrawable) isCaretPrintable(text text.Line, textIndex int, part string, partIndex int) bool {
	fragment := text.Text[textIndex]

	isCaret := len(part) == 0 && fragment.Atom.HasAny(style.AtmSelect)
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

func (d *TextAreaDrawable) draw() ([]text.Line, bool) {
	assert.True(d.initialized, "the drawable should be initialized before draw")

	return d.drawable.Draw()
}
