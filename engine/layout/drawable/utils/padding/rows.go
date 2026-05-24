package padding

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding/options"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type rowPositioner func([]text.Line, text.Fragment, winsize.Rows) []text.Line

var rowPositionerMap = map[style.VerticalPosition]rowPositioner{
	style.Top:    rowsToTop,
	style.Bottom: rowsToBottom,
	style.Middle: rowsToMiddle,
}

func Rows(rows hint.Size[winsize.Rows], opts ...options.RowsOption) transformer {
	cfg := options.ResolveRowsConfig(opts...)

	return func(size winsize.Winsize, lines []text.Line) []text.Line {
		frag := cfg.Provider(size, lines...)

		padding := rows.Min(size.Rows)
		if winsize.Rows(len(lines)) >= padding {
			return lines
		}

		positioner, ok := rowPositionerMap[cfg.Position]
		if !ok {
			assert.Unreachable("unhandled vertical position '%d'", cfg.Position)
			positioner = rowsToTop
		}

		return positioner(lines, frag, padding)
	}
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
