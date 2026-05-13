package wrap

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func NormalizeLines(lines ...text.Line) []LayoutLine {
	return normalizeLines(false, lines...)
}

func NormalizeLinesWithOrder(lines ...text.Line) []LayoutLine {
	return normalizeLines(true, lines...)
}

func normalizeLines(order bool, lines ...text.Line) []LayoutLine {
	buffer := make([]LayoutLine, 0, len(lines))

	for _, line := range lines {
		normalizedLF := splitLineFeeds(&line, order)

		for _, n := range normalizedLF {
			words := splitLineWords(&n)
			newLayoutLine := NewLayoutLine(&n, words...)
			buffer = append(buffer, *newLayoutLine)
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

func splitLineFeeds(line *text.Line, order bool) []text.Line {
	result := make([]text.Line, 0)

	index := uint16(1)
	if line.Order != 0 {
		index = line.Order
	}

	current := text.LineFromMeta(line)
	if order {
		current.SetOrder(index)
	}

	for _, frag := range line.Text {
		if !strings.ContainsAny(frag.Text, "\n\r") {
			current.PushFragments(frag)
			continue
		}

		normalizedText := runes.NormalizeLineFeed(frag.Text)

		parts := strings.Split(normalizedText, "\n")
		for i, part := range parts {
			if part != "" {
				current.PushFragments(
					*text.NewFragment(part).CopyMeta(&frag),
				)
			}

			if i >= len(parts)-1 {
				continue
			}

			result = append(result, *current)
			index += 1

			current = text.LineFromMeta(line)
			if order {
				current.SetOrder(index)
			}
		}
	}

	result = append(result, *current)

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
