package line

import (
	"github.com/Rafael24595/go-terminal/engine/core"
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

func indexLines(cols int, line core.Line, meta *IndexMeta) []core.Line {
	isGreaterWithoutIndex := line.Len() > int(cols)
	isGreaterWithIndex := meta != nil && line.Len()+int(meta.totalWidth) > cols

	if isGreaterWithoutIndex || isGreaterWithIndex {
		return WrapLineWordsWithIndex(int(cols), line, meta)
	}

	fragments := core.FragmentsFromString()
	if meta != nil {
		fragments = append(fragments, core.NewFragment(meta.header(int(line.Order))))
	}

	newLine := core.LineFromFragments(
		append(fragments, line.Text...)...,
	)

	return core.FixedLinesFromLines(line.Spec, newLine)
}

func WrapLineWords(cols int, line core.Line) []core.Line {
	return WrapLineWordsWithIndex(cols, line, nil)
}

func WrapLineWordsWithIndex(cols int, line core.Line, meta *IndexMeta) []core.Line {
	result := make([]core.Line, 0)
	current := core.LineFromSpec(line.Spec)
	width := 0

	words := core.TokenizeLineWords(line)

	if meta != nil {
		fragments := core.FragmentsFromString(meta.header(int(line.Order)))
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
			current = core.LineFromSpec(line.Spec)

			if meta != nil {
				fragments := core.FragmentsFromString(meta.body())
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
	word core.WordToken,
	cols int,
	current core.Line,
	width int,
	meta *IndexMeta,
) (core.Line, []core.Line, int) {
	current, lines, width := core.SplitLongToken(word, cols, current, width)
	if meta == nil || len(lines) == 0 {
		return current, lines, width
	}

	index := core.FragmentsFromString(meta.body())

	current.Text = append(index, current.Text...)

	for i := 1; i < len(lines); i++ {
		lines[i].Text = append(index, lines[i].Text...)
	}

	return current, lines, width
}

func computeIndexMeta(lines []core.Line) *IndexMeta {
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
		prefixBody: helper.FillRight(" ", int(size)),
		digits:     uint16(size),
		totalWidth: size + uint32(len(separator)),
	}
}
