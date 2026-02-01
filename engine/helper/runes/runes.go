package runes

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

type RuneDefinition struct {
	Rune rune
	Skip bool
}

func NormalizeLineEnd(text string) string {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(normalized, "\r", "\n")
}

func AppendAt(slice []rune, insert []rune, pos uint) []rune {
	i := int(pos)
	size := len(insert)
	slice = append(slice, make([]rune, size)...)
	copy(slice[i+size:], slice[i:])
	copy(slice[i:], insert)
	return slice
}

func AppendRange(slice []rune, insert []rune, start, end uint) []rune {
	if start == end {
		return AppendAt(slice, insert, start)
	}

	s := int(start)
	e := int(end)

	oldSize := len(slice)
	newSize := len(insert)

	newSlice := make([]rune, oldSize+(newSize-e+s))

	copy(newSlice[0:s], slice[0:s])
	copy(newSlice[s:], insert)
	copy(newSlice[s+newSize:], slice[e:])

	return newSlice
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
