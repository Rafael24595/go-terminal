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

func Line(cols winsize.Cols, line *text.Line) []text.Line {
	if cols >= text.FragmentMeasure(cols, line.Text...) {
		return []text.Line{*line}
	}

	result := make([]text.Line, 0)
	current := text.EmptyLine().AddSpec(line.Spec)
	width := winsize.Cols(0)

	words := splitLineWords(line)

	for _, word := range words {
		wordlen := text.FragmentMeasure(cols, word.Text...)

		if width+wordlen <= cols {
			current.Text = append(current.Text, word.Text...)
			width += wordlen

			continue
		}

		if wordlen <= cols {
			result = append(result, *current)
			current = text.EmptyLine().
				AddSpec(line.Spec)

			current.Text = append(current.Text, word.Text...)
			width = wordlen

			continue
		}

		newCurrent, lines, newWidth := splitLongWord(word, cols, *current, width)

		result = append(result, lines...)
		current = &newCurrent
		width = newWidth
	}

	if len(current.Text) > 0 {
		result = append(result, *current)
	}

	return result
}

func NextLine(cols winsize.Cols, lines []text.Line) (*text.Line, []text.Line) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]text.Line, 0)
	}

	target := lines[0]
	remain := lines[1:]

	cursor := text.EmptyLine().
		CopyMeta(&target)

	remaining := cols

	for len(target.Text) > 0 {
		frag := target.Text[0]
		fragMeasure := text.FragmentMeasure(cols, frag)

		if fragMeasure <= remaining {
			cursor.PushFragments(frag)
			remaining = remaining.Clamp(fragMeasure)
			target.Text = target.Text[1:]
			continue
		}

		if len(cursor.Text) == 0 && remaining > 0 {
			taken, restFrag := splitFragmentAt(&frag, remaining)
			cursor.PushFragments(*taken)
			target.Text[0] = *restFrag
		}

		newRest := append([]text.Line{target}, remain...)
		return cursor, newRest
	}

	return cursor, remain
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
