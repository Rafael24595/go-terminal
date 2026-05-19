package padding

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type colPositioner func(winsize.Cols) (winsize.Cols, winsize.Cols)

var colPositionerMap = map[style.HorizontalPosition]colPositioner{
	style.Left:   colToLeft,
	style.Right:  colToRight,
	style.Center: colToCenter,
}

func colToLeft(remaining winsize.Cols) (winsize.Cols, winsize.Cols) {
	return 0, remaining
}

func colToRight(remaining winsize.Cols) (winsize.Cols, winsize.Cols) {
	return remaining, 0
}

func colToCenter(remaining winsize.Cols) (winsize.Cols, winsize.Cols) {
	paddingL := remaining / 2
	paddingR := remaining.Sub(paddingL)
	return paddingL, paddingR
}

func Cols(cols SizeHint[winsize.Cols], position ...style.HorizontalPosition) transformer {
	horizontal := style.Left
	if len(position) > 0 {
		horizontal = position[0]
	}

	return func(size winsize.Winsize, lines []text.Line) []text.Line {
		newLines := make([]text.Line, len(lines))
		fixedMin := cols.min(size.Cols)

		for i := range lines {
			remaining := fixedMin.Sub(
				text.FragmentMeasure(size.Cols, lines[i].Text...),
			)

			if remaining == 0 {
				newLines[i] = lines[i]
				continue
			}

			newLines[i] = addColsPadding(remaining, lines[i], horizontal)
		}

		return newLines
	}
}

func addColsPadding(
	cols winsize.Cols,
	line text.Line,
	position style.HorizontalPosition,
) text.Line {
	positioner, ok := colPositionerMap[position]
	if !ok {
		assert.Unreachable("undefined vertical position '%d'", position)
		positioner = colToLeft
	}

	paddingL, paddingR := positioner(cols)

	frags := make([]text.Fragment, 0, 3)

	if paddingL > 0 {
		frag := text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(paddingL))
		frags = append(frags, *frag)
	}

	frags = append(frags, line.Text...)

	if paddingR > 0 {
		frag := text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(paddingR))
		frags = append(frags, *frag)
	}

	return *text.LineFromMeta(&line).
		PushFragments(frags...)
}
