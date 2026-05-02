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

func AppendAt(slice []rune, insert []rune, position offset.Offset) []rune {
	size := offset.Offset(len(insert))
	slice = append(slice, make([]rune, size)...)
	copy(slice[position+size:], slice[position:])
	copy(slice[position:], insert)
	return slice
}

func AppendRange(slice []rune, insert []rune, start, end offset.Offset) []rune {
	if start == end {
		return AppendAt(slice, insert, start)
	}

	sliceLen := offset.Offset(len(slice))
	insertLen := offset.Offset(len(insert))

	assert.False(
		sliceLen < end,
		"range[%d - %d] is greater than slice length %d", start, end, sliceLen,
	)

	size := sliceLen.Clamp(
		end.Clamp(start),
	)

	newSlice := make([]rune, size+insertLen)

	copy(newSlice[0:start], slice[0:start])
	copy(newSlice[start:], insert)
	if end < sliceLen {
		copy(newSlice[start+insertLen:], slice[end:])
	}

	return newSlice
}

func NormalizeBuffer(buf []rune, min uint) []rune {
	if uint(len(buf)) >= min {
		return buf
	}

	needed := int(min) - len(buf)
	padding := make([]rune, needed)

	return append(buf, padding...)
}

func BackwardIndexWithLimit[T math.Number](b []rune, rs []RuneDefinition, i T) T {
	return BackwardIndex(b, rs, i) + 1
}

func BackwardIndex[T math.Number](b []rune, rs []RuneDefinition, i T) T {
	s := fixdBackwardIndex(b, rs, i)

	for j := s - 1; j >= 0; j-- {
		for _, v := range rs {
			if v.Rune == b[j] {
				return T(j) + 1
			}
		}
	}

	return 0
}

func fixdBackwardIndex[T math.Number](b []rune, rs []RuneDefinition, i T) int {
	s := max(0, int(i)-1)

	for s > 0 {
		for _, v := range rs {
			if v.Rune != b[s] || !v.Skip {
				return s
			}
		}

		s--
	}

	return s
}

func ForwardIndexWithLimit[T math.Number](b []rune, rs []RuneDefinition, i T) T {
	s := i

	if s < T(len(b)) {
		for _, v := range rs {
			if b[s] == v.Rune {
				return s
			}
		}
	}

	return ForwardIndex(b, rs, s)
}

func ForwardIndex[T math.Number](b []rune, rs []RuneDefinition, i T) T {
	s := fixForwardIndex(b, rs, i)

	for s < len(b) {
		for _, v := range rs {
			if v.Rune == b[s] {
				return T(s)
			}
		}

		s++
	}

	return T(s)
}

func fixForwardIndex[T math.Number](b []rune, rs []RuneDefinition, i T) int {
	s := min(int(i+1), len(b))

	for s < len(b) {
		for _, v := range rs {
			if v.Rune != b[s] || !v.Skip {
				return s
			}
		}

		s++
	}

	return s
}

func JoinReverse(ps []string) string {
	var sb strings.Builder
	for i := len(ps) - 1; i >= 0; i-- {
		sb.WriteString(ps[i])
	}
	return sb.String()
}

func RuneIndexToByteIndex(text string, runeIndex offset.Offset) (offset.Offset, bool) {
	if runeIndex == 0 {
		return 0, true
	}

	count := offset.Offset(0)
	for i := range text {
		if count == runeIndex {
			return offset.Offset(i), true
		}
		count++
	}

	if count == runeIndex {
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
