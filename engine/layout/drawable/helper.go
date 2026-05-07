package drawable

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func MaxLineSize(cols winsize.Cols, lines ...text.Line) winsize.Cols {
	size := winsize.Cols(0)
	for _, l := range lines {
		measure := text.FragmentMeasure(cols, l.Text...)
		size = max(size, measure)
	}
	return size
}

func DrainDrawable(
	size winsize.Winsize,
	drawable Drawable,
	lazy bool,
) ([]text.Line, bool) {
	result := make([]text.Line, 0, size.Rows)

	remaining := size.Rows
	for {
		lines, hasNext := drawable.Draw(size)
		linesLen := winsize.Rows(len(lines))

		limit := linesLen
		if lazy {
			limit = min(linesLen, remaining)
		}

		result = append(result, lines[:limit]...)

		remaining = math.SubClampZero(
			remaining,
			winsize.Rows(len(lines)),
		)

		shouldExit := lazy && remaining == 0
		if !hasNext || shouldExit {
			hasRest := limit < linesLen
			return result, hasNext || hasRest
		}
	}
}
