package wrap

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func NormalizeLines(lines ...text.Line) []text.Line {
	buffer := make([]text.Line, 0, len(lines))

	for _, l := range lines {
		normalizedLF := splitLineFeeds(&l)

		for _, n := range normalizedLF {
			newLine := text.EmptyLine().
				CopyMeta(&n)

			for _, w := range splitLineWords(&n) {
				newLine.PushFragments(w.Text...)
			}

			buffer = append(buffer, *newLine)
		}
	}

	return buffer
}

func MaterializeEmpty(
	size winsize.Winsize,
	placeholder string,
	lines ...LayoutLine,
) []LayoutLine {
	for i, line := range lines {
		if text.FragmentMeasure(size.Cols, line.Source.Text...) != 0 {
			continue
		}

		lastFrag := text.Fragment{}
		if len(line.Source.Text) > 0 {
			lastFrag = line.Source.Text[len(line.Source.Text)-1]
		}

		fragment := *text.NewFragment(placeholder).
			CopyMeta(&lastFrag)

		lines[i].Source.PushFragments(fragment)
		lines[i].Words = append(line.Words, *newWord(fragment))
	}

	return lines
}

func Line(cols winsize.Cols, line *text.Line) []text.Line {
	result := make([]text.Line, 0)
	current := line

	for current != nil {
		head, rest := wrapOnceFromLine(cols, *current)
		result = append(result, *head)
		current = rest
	}

	return result
}

func NextLine(cols winsize.Cols, lines []LayoutLine) (*text.Line, []LayoutLine) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]LayoutLine, 0)
	}

	current := lines[0]
	remain := lines[1:]

	result, rest := wrapOnce(cols, current)
	if rest != nil {
		remain = append([]LayoutLine{*rest}, remain...)
	}

	return result, remain
}

func wrapOnceFromLine(cols winsize.Cols, line text.Line) (*text.Line, *text.Line) {
	words := splitLineWords(&line)

	layout := NewLayoutLine(&line, words...)

	result, rest := wrapOnce(cols, *layout)
	if rest == nil {
		return result, nil
	}

	return result, rest.toLine()
}

func wrapOnce(cols winsize.Cols, line LayoutLine) (*text.Line, *LayoutLine) {
	cursor := text.LineFromMeta(line.Source)

	remaining := cols

	words := line.Words

	for len(words) > 0 {
		word := words[0]
		wordMeasure := text.FragmentMeasure(cols, word.Text...)

		if wordMeasure <= remaining {
			cursor.PushFragments(word.Text...)
			remaining = remaining.Clamp(wordMeasure)
			words = words[1:]

			continue
		}

		if text.FragmentMeasure(cols, cursor.Text...) > 0 {
			break
		}

		newWord, restWord := splitLongWord(word, cols, remaining)
		if newWord != nil {
			cursor.PushFragments(newWord.Text...)
		}

		var rest *LayoutLine
		if restWord != nil {
			rest = NewLayoutLine(line.Source, *restWord)
		}

		return cursor, rest
	}

	if len(words) == 0 {
		return cursor, nil
	}

	rest := NewLayoutLine(line.Source, words...)

	return cursor, rest
}

func splitLineFeeds(line *text.Line) []text.Line {
	result := make([]text.Line, 0)

	current := text.EmptyLine().
		CopyMeta(line)

	for _, frag := range line.Text {
		if !strings.Contains(frag.Text, "\n") && !strings.Contains(frag.Text, "\r") {
			current.PushFragments(frag)
			continue
		}

		normalizedText := runes.NormalizeLineFeed(frag.Text)

		parts := strings.Split(normalizedText, "\n")
		for i, part := range parts {
			newFrag := text.NewFragment(part).
				CopyMeta(&frag)

			current.PushFragments(*newFrag)

			if i < len(parts)-1 {
				result = append(result, *current)

				current = text.EmptyLine().
					CopyMeta(line)
			}
		}
	}

	if len(current.Text) > 0 {
		result = append(result, *current)
	}

	return result
}

func splitFragmentAt(frag *text.Fragment, cols winsize.Cols) (*text.Fragment, *text.Fragment) {
	if cols <= 0 {
		return text.EmptyFragment().
			CopyMeta(frag), frag
	}

	byteIndex, canBreak := runes.RuneIndexToByteIndex(frag.Text, offset.Offset(cols))
	if !canBreak {
		return frag, text.EmptyFragment().
			CopyMeta(frag)
	}

	taken := text.NewFragment(frag.Text[:byteIndex]).
		CopyMeta(frag)

	rest := text.NewFragment(frag.Text[byteIndex:]).
		CopyMeta(frag)

	return taken, rest
}
