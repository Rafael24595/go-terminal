package padding

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type colPositioner func(winsize.Cols) (winsize.Cols, winsize.Cols)

var colPositionerMap = map[style.HorizontalPosition]colPositioner{
	style.Left:   colToLeft,
	style.Right:  colToRight,
	style.Center: colToCenter,
}

func Cols(cols hint.Size[winsize.Cols], opts ...cols.Option) transformer {
	return func(size winsize.Winsize, lines []text.Line) []text.Line {
		newLines := make([]text.Line, len(lines))
		fixedMin := cols.Min(size.Cols)

		for i := range lines {
			remaining := fixedMin.Sub(
				text.FragmentMeasure(size.Cols, lines[i].Text...),
			)

			if remaining == 0 {
				newLines[i] = lines[i]
				continue
			}

			newLines[i] = AddColsPadding(remaining, lines[i], opts...)
		}

		return newLines
	}
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

func AddColsPadding(
	size winsize.Cols,
	line text.Line,
	opts ...cols.Option,
) text.Line {
	cfg := cols.ResolveConfig(opts...)

	frag := cfg.Provider(size, line)

	positioner, ok := colPositionerMap[cfg.Position]
	if !ok {
		assert.Unreachable("undefined vertical position '%d'", cfg.Position)
		positioner = colToLeft
	}

	paddingL, paddingR := positioner(size)

	frags := make([]text.Fragment, 0, 3)

	if paddingL > 0 {
		frag := frag.Clone().
			AddSpec(style.SpecRepeatRight(paddingL))
		frags = append(frags, *frag)
	}

	frags = append(frags, line.Text...)

	if paddingR > 0 {
		frag := frag.Clone().
			AddSpec(style.SpecRepeatRight(paddingR))
		frags = append(frags, *frag)
	}

	return *text.LineFromMeta(&line).
		PushFragments(frags...)
}
