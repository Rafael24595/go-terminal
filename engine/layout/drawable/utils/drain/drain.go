package drain

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/work"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

func DrawableEager(size winsize.Winsize, drawable drawable.Drawable) []text.Line {
	result, _ := Drawable(size, drawable, false)
	return result
}

func DrawableLazy(size winsize.Winsize, drawable drawable.Drawable) ([]text.Line, bool) {
	return Drawable(size, drawable, true)
}

func Drawable(
	size winsize.Winsize,
	drawable drawable.Drawable,
	lazy bool,
) ([]text.Line, bool) {
	result := make([]text.Line, 0, size.Rows)
	if size.Rows == 0 {
		return result, false
	}

	tracker := work.NewTracker()
	tracker.Add(1)

	remaining := uint(size.Rows)
	for tracker.Unfinished() {
		tracker.Advance()
		tracker.Reset()

		lines, hasNext := drawable.Draw(size)
		if hasNext {
			tracker.Add(1)
		}

		linesLen := uint(len(lines))
		if linesLen == 0 {
			return result, false
		}

		tracker.Add(linesLen)

		for _, lne := range lines {
			fixed := wrap.Line(size.Cols, &lne)

			fixedLen := uint(len(fixed))
			if fixedLen == 0 {
				continue
			}

			tracker.Advance()
			tracker.Add(fixedLen)

			for _, fix := range fixed {
				tracker.Advance()

				result = append(result, fix)

				if !lazy {
					continue
				}

				remaining = math.SubClampZero(remaining, 1)
				if remaining == 0 {
					return result, tracker.Unfinished()
				}
			}
		}
	}

	return result, false
}
