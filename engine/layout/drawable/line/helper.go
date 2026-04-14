package line

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/text"
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

func TokenizeLines(lines ...text.Line) []text.Line {
	buffer := make([]text.Line, len(lines))

	for i, l := range lines {
		newLine := text.EmptyLine().CopyMeta(&l)
		for _, w := range text.TokenizeLineWords(&l) {
			newLine.PushFragments(w.Text...)
		}
		buffer[i] = *newLine
	}

	return buffer
}

func WrapLineWords(cols int, line *text.Line) []text.Line {
	if cols >= text.FragmentMeasure(line.Text...) {
		return []text.Line{*line}
	}

	result := make([]text.Line, 0)
	current := text.EmptyLine().AddSpec(line.Spec)
	width := 0

	words := text.TokenizeLineWords(line)

	for _, word := range words {
		wordlen := word.Size()

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

		newCurrent, lines, newWidth := text.SplitLongToken(word, cols, *current, width)

		result = append(result, lines...)
		current = &newCurrent
		width = newWidth
	}

	if len(current.Text) > 0 {
		result = append(result, *current)
	}

	return result
}

func WrapNextLine(cols uint16, lines []text.Line, meta *IndexMeta) (*text.Line, []text.Line) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]text.Line, 0)
	}

	target := lines[0]
	remain := lines[1:]

	cursor := text.EmptyLine().CopyMeta(&target)

	width := int(cols)

	emptyLen := 0
	if meta != nil {
		var prefix string
		if target.Order != 0 {
			prefix = meta.header(int(target.Order))
			target.Order = 0
		} else {
			prefix = meta.body()
		}

		cursor.PushFragments(*text.NewFragment(prefix))
		width = math.SubClampZero(width, int(meta.totalWidth))

		emptyLen = len(cursor.Text)
	}

	for len(target.Text) > 0 {
		frag := target.Text[0]
		fragMeasure := text.FragmentMeasure(frag)

		if fragMeasure <= width {
			cursor.PushFragments(frag)
			width = math.SubClampZero(width, fragMeasure)
			target.Text = target.Text[1:]
			continue
		}

		if len(cursor.Text) == emptyLen && width > 0 {
			taken, restFrag := text.TakeFromFragment(&frag, width)
			cursor.PushFragments(*taken)
			target.Text[0] = *restFrag

			newRest := append([]text.Line{target}, remain...)
			return cursor, newRest
		}

		if len(cursor.Text) == 1 && meta != nil {
			assert.Unreachable("index prefix should be lesser than line size")

			cursor.PushFragments(frag)
			target.Text = target.Text[1:]
		}

		newRest := append([]text.Line{target}, remain...)
		return cursor, newRest
	}

	return cursor, remain
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
