package padding

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func DefaultRowFrag() *text.Fragment {
	return text.EmptyFragment()
}

type rowPositioner func([]text.Line, text.Fragment, winsize.Rows) []text.Line

var rowPositionerMap = map[style.VerticalPosition]rowPositioner{
	style.Top:    rowsToTop,
	style.Bottom: rowsToBottom,
	style.Middle: rowsToMiddle,
}

func rowsToTop(lines []text.Line, frag text.Fragment, padding winsize.Rows) []text.Line {
	newLines := paddingLines(padding, frag)
	copy(newLines, lines)
	return newLines
}

func rowsToBottom(lines []text.Line, frag text.Fragment, padding winsize.Rows) []text.Line {
	rest := padding.Sub(
		winsize.Rows(len(lines)),
	)

	newLines := paddingLines(rest, frag)
	return append(newLines, lines...)
}

func rowsToMiddle(lines []text.Line, frag text.Fragment, padding winsize.Rows) []text.Line {
	rest := padding.Sub(
		winsize.Rows(len(lines)),
	)

	half := rest / 2

	top := paddingLines(half, frag)
	bottom := paddingLines(rest.Sub(half), frag)

	newLines := append(top, lines...)
	return append(newLines, bottom...)
}

func paddingLines(rows winsize.Rows, frag text.Fragment) []text.Line {
	result := make([]text.Line, rows)

	for i := range result {
		result[i].Text = append(result[i].Text, frag)
	}

	return result
}

func Rows(rows hint.Size[winsize.Rows], position ...style.VerticalPosition) transformer {
	return RowsWithDefault(rows, *DefaultRowFrag(), position...)
}

func RowsWithDefault(rows hint.Size[winsize.Rows], frag text.Fragment, position ...style.VerticalPosition) transformer {
	vertical := style.Top
	if len(position) > 0 {
		vertical = position[0]
	}

	return func(size winsize.Winsize, lines []text.Line) []text.Line {
		padding := rows.Min(size.Rows)
		if winsize.Rows(len(lines)) >= padding {
			return lines
		}

		positioner, ok := rowPositionerMap[vertical]
		if !ok {
			assert.Unreachable("unhandled vertical position '%d'", position)
			positioner = rowsToTop
		}

		return positioner(lines, frag, padding)
	}
}
