package line

import (
	"github.com/Rafael24595/go-terminal/engine/core/marker"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

const separator = " | "

type IndexMeta struct {
	sufix      string
	prefixBody string
	digits     uint16
	totalWidth uint32
}

func (i IndexMeta) header(index int) string {
	return helper.Right(index, int(i.digits)) + i.sufix
}

func (i IndexMeta) body() string {
	return i.prefixBody + i.sufix
}

func indexLines(cols int, line text.Line, meta *IndexMeta) []text.Line {
	measure := text.LineFragmentsMeasure(line)

	isGreaterWithoutIndex := measure > int(cols)
	isGreaterWithIndex := meta != nil && measure+int(meta.totalWidth) > cols

	if isGreaterWithoutIndex || isGreaterWithIndex {
		return WrapLineWordsWithIndex(int(cols), line, meta)
	}

	fragments := text.FragmentsFromString()
	if meta != nil {
		fragments = append(fragments, text.NewFragment(meta.header(int(line.Order))))
	}

	newLine := text.LineFromFragments(
		append(fragments, line.Text...)...,
	)

	return text.FixedLinesFromLines(line.Spec, newLine)
}

func WrapLineWords(cols int, line text.Line) []text.Line {
	if cols >= text.LineFragmentsMeasure(line) {
		return []text.Line{line}
	}
	return WrapLineWordsWithIndex(cols, line, nil)
}

func WrapLineWordsWithIndex(cols int, line text.Line, meta *IndexMeta) []text.Line {
	result := make([]text.Line, 0)
	current := text.LineFromSpec(line.Spec)
	width := 0

	words := text.TokenizeLineWords(line)

	if meta != nil {
		fragments := text.FragmentsFromString(meta.header(int(line.Order)))
		current.Text = append(current.Text, fragments...)
		cols -= int(meta.totalWidth)
	}

	for _, word := range words {
		wordlen := word.Size()

		if width+wordlen <= cols {
			current.Text = append(current.Text, word.Text...)
			width += wordlen

			continue
		}

		if wordlen <= cols {
			result = append(result, current)
			current = text.LineFromSpec(line.Spec)

			if meta != nil {
				fragments := text.FragmentsFromString(meta.body())
				current.Text = append(current.Text, fragments...)
			}

			current.Text = append(current.Text, word.Text...)
			width = wordlen

			continue
		}

		newCurrent, lines, newWidth := wrapLongTokenWithIndex(word, cols, current, width, meta)

		result = append(result, lines...)
		current = newCurrent
		width = newWidth
	}

	if len(current.Text) > 0 {
		result = append(result, current)
	}

	return result
}

func wrapLongTokenWithIndex(
	word text.WordToken,
	cols int,
	current text.Line,
	width int,
	meta *IndexMeta,
) (text.Line, []text.Line, int) {
	current, lines, width := text.SplitLongToken(word, cols, current, width)
	if meta == nil || len(lines) == 0 {
		return current, lines, width
	}

	index := text.FragmentsFromString(meta.body())

	current.Text = append(index, current.Text...)

	for i := 1; i < len(lines); i++ {
		lines[i].Text = append(index, lines[i].Text...)
	}

	return current, lines, width
}

func computeIndexMeta(lines []text.Line) *IndexMeta {
	size := uint32(0)

	for _, line := range lines {
		if line.Order == 0 {
			continue
		}
		size = max(size, math.Digits(line.Order))
	}

	if size == 0 {
		return nil
	}

	return &IndexMeta{
		sufix:      separator,
		prefixBody: helper.FillRight(marker.DefaultPaddingText, int(size)),
		digits:     uint16(size),
		totalWidth: size + uint32(len(separator)),
	}
}
