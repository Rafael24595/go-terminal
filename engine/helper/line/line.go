package line

import "github.com/Rafael24595/go-terminal/engine/core/key"

func DistanceFromLF(buffer []rune, from int) int {
	return from - FindLineStart(buffer, from)
}

func FindLineStart(buf []rune, from int) int {
	for i := from - 1; i >= 0; i-- {
		if buf[i] == key.ENTER_LF {
			return i + 1
		}
	}
	return 0
}

func FindLineEnd(buf []rune, start int) int {
	i := start
	for i < len(buf) && buf[i] != key.ENTER_LF {
		i++
	}
	return i
}

func FindNextLineStart(buf []rune, from int) int {
	for i := from; i < len(buf); i++ {
		if buf[i] == key.ENTER_LF {
			return i + 1
		}
	}
	return -1
}

func FindPrevLineStart(buf []rune, from int) int {
	prevLineStart := FindLineStart(buf, from)
	if prevLineStart == 0 {
		return -1
	}
	return FindLineStart(buf, prevLineStart-1)
}

func ClampToLine(buf []rune, lineStart, col int) int {
	end := FindLineEnd(buf, lineStart)
	lineLen := end - lineStart

	if col > lineLen {
		return end
	}

	return lineStart + col
}
