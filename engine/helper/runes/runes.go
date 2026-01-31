package runes

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
)

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

func BackwardIndex[T math.Number](b []rune, r rune, i T) T {
	s := max(0, int(i)-1)
	for s > 0 && b[s-1] == r {
		s--
	}

	for i := s - 1; i >= 0; i-- {
		if b[i] == r {
			return T(i) + 1
		}
	}

	return 0
}

func ForwardIndex[T math.Number](b []rune, r rune, i T) T {
	s := min(int(i+1), len(b))

	for s < len(b) && b[s] == r {
		s++
	}

	newEnd := s
	for newEnd < len(b) && b[newEnd] != r {
		newEnd++
	}

	return T(newEnd)
}
