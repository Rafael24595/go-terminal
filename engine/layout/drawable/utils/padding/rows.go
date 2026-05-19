package padding

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type rowPositioner func([]text.Line, winsize.Rows) []text.Line

var rowPositionerMap = map[style.VerticalPosition]rowPositioner{
	style.Top:    rowsToTop,
	style.Bottom: rowsToBottom,
	style.Middle: rowsToMiddle,
}

func rowsToTop(lines []text.Line, padding winsize.Rows) []text.Line {
	newLines := make([]text.Line, padding)
	copy(newLines, lines)
	return newLines
}

func rowsToBottom(lines []text.Line, padding winsize.Rows) []text.Line {
	rest := padding.Sub(
		winsize.Rows(len(lines)),
	)

	newLines := make([]text.Line, rest)
	return append(newLines, lines...)
}

func rowsToMiddle(lines []text.Line, padding winsize.Rows) []text.Line {
	rest := padding.Sub(
		winsize.Rows(len(lines)),
	)

	half := rest / 2

	top := make([]text.Line, half)
	bottom := make([]text.Line, rest.Sub(half))

	newLines := append(top, lines...)
	return append(newLines, bottom...)
}

func Rows(rows SizeHint[winsize.Rows], position ...style.VerticalPosition) transformer {
	vertical := style.Top
	if len(position) > 0 {
		vertical = position[0]
	}

	return func(size winsize.Winsize, lines []text.Line) []text.Line {
		padding := rows.min(size.Rows)
		if winsize.Rows(len(lines)) >= padding {
			return lines
		}

		positioner, ok := rowPositionerMap[vertical]
		if !ok {
			assert.Unreachable("unhandled vertical position '%d'", position)
			positioner = rowsToTop
		}

		return positioner(lines, padding)
	}
}
