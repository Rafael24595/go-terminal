package runes

import (
	"strings"
	"unicode/utf8"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

var NextWordRunes = []RuneDefinition{
	{
		Rune: ' ',
		Skip: false,
	},
	{
		Rune: '.',
		Skip: true,
	},
	{
		Rune: ',',
		Skip: true,
	},
	{
		Rune: ascii.ENTER_LF,
		Skip: true,
	},
}

var NextLineRunes = []RuneDefinition{
	{
		Rune: ascii.ENTER_LF,
		Skip: false,
	},
}

type RuneDefinition struct {
	Rune rune
	Skip bool
}

func NormalizeLineEnd(text string) string {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(normalized, "\r", "\n")
}

func AppendAt(
	buffer []rune,
	insert []rune,
	position offset.Offset,
) []rune {
	size := offset.Offset(len(insert))
	buffer = append(buffer, make([]rune, size)...)
	copy(buffer[position+size:], buffer[position:])
	copy(buffer[position:], insert)
	return buffer
}

func AppendRange(
	buffer []rune,
	insert []rune,
	start, end offset.Offset,
) []rune {
	if start == end {
		return AppendAt(buffer, insert, start)
	}

	sliceLen := offset.Offset(len(buffer))
	insertLen := offset.Offset(len(insert))

	assert.False(
		sliceLen < end,
		"range[%d - %d] is greater than slice length %d", start, end, sliceLen,
	)

	size := sliceLen.Clamp(
		end.Clamp(start),
	)

	newSlice := make([]rune, size+insertLen)

	copy(newSlice[0:start], buffer[0:start])
	copy(newSlice[start:], insert)
	if end < sliceLen {
		copy(newSlice[start+insertLen:], buffer[end:])
	}

	return newSlice
}

func NormalizeBuffer(buffer []rune, min uint) []rune {
	bufferLen := uint(len(buffer))
	if bufferLen >= min {
		return buffer
	}

	needed := math.SubClampZero(min, bufferLen)
	padding := make([]rune, needed)

	return append(buffer, padding...)
}

func BackwardIndexWithLimit[T math.Number](
	buffer []rune,
	definition []RuneDefinition,
	index T,
) T {
	return BackwardIndex(buffer, definition, index) + 1
}

func BackwardIndex[T math.Number](
	buffer []rune,
	definition []RuneDefinition,
	index T,
) T {
	newIndex := fixdBackwardIndex(buffer, definition, index)

	for j := newIndex - 1; j >= 0; j-- {
		for _, v := range definition {
			if v.Rune == buffer[j] {
				return j + 1
			}
		}
	}

	return 0
}

func fixdBackwardIndex[T math.Number](
	buffer []rune,
	definition []RuneDefinition,
	index T,
) T {
	newIndex := math.SubClampZero(index, 1)

	for newIndex > 0 {
		for _, v := range definition {
			if v.Rune != buffer[newIndex] || !v.Skip {
				return newIndex
			}
		}

		newIndex--
	}

	return newIndex
}

func ForwardIndexWithLimit[T math.Number](
	buffer []rune,
	definition []RuneDefinition,
	index T,
) T {
	if index < T(len(buffer)) {
		for _, v := range definition {
			if buffer[index] == v.Rune {
				return index
			}
		}
	}

	return ForwardIndex(buffer, definition, index)
}

func ForwardIndex[T math.Number](
	buffer []rune,
	definition []RuneDefinition,
	index T,
) T {
	newIndex := fixForwardIndex(buffer, definition, index)

	for newIndex < T(len(buffer)) {
		for _, v := range definition {
			if v.Rune == buffer[newIndex] {
				return T(newIndex)
			}
		}

		newIndex++
	}

	return newIndex
}

func fixForwardIndex[T math.Number](
	buffer []rune,
	definition []RuneDefinition,
	index T,
) T {
	bufferLen := T(len(buffer))
	newIndex := min(index+1, bufferLen)

	for newIndex < bufferLen {
		for _, v := range definition {
			if v.Rune != buffer[newIndex] || !v.Skip {
				return newIndex
			}
		}

		newIndex++
	}

	return newIndex
}

func JoinReverse(buffer []string) string {
	var sb strings.Builder
	for i := len(buffer) - 1; i >= 0; i-- {
		sb.WriteString(buffer[i])
	}
	return sb.String()
}

func RuneIndexToByteIndex(text string, index offset.Offset) (offset.Offset, bool) {
	if index == 0 {
		return 0, true
	}

	count := offset.Offset(0)
	for i := range text {
		if count == index {
			return offset.Offset(i), true
		}
		count++
	}

	if count == index {
		return offset.Offset(len(text)), true
	}

	return 0, false

}

func Measure(text string) winsize.Cols {
	return winsize.Cols(measure(text))
}

func Measureo(text string) offset.Offset {
	return offset.Offset(measure(text))
}

func measure(text string) int {
	return utf8.RuneCountInString(text)
}
