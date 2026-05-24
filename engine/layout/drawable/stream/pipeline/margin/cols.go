package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func ColsLeft(margin winsize.Cols) pipeline.DataTransformer {
	return colsDrawTransformer(margin, cols.WithPosition(style.Right))
}

func ColsRight(margin winsize.Cols) pipeline.DataTransformer {
	return colsDrawTransformer(margin, cols.WithPosition(style.Left))
}

func ColsCenter(margin winsize.Cols) pipeline.DataTransformer {
	return colsDrawTransformer(margin, cols.WithPosition(style.Center))
}

func colsDrawTransformer(margin winsize.Cols, opts ...cols.Option) pipeline.DataTransformer {
	cfg := cols.ResolveConfig(opts...)

	margin = margin * horizontalFactor(cfg.Position)

	return func(size winsize.Winsize, _ drawable.Unit, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		newLines := make([]text.Line, len(lines))

		for i := range lines {
			measure := text.FragmentMeasure(size.Cols, lines[i].Text...) + margin

			cols := size.Cols + margin
			if cols.Sub(measure) == 0 {
				newLines[i] = lines[i]
				continue
			}

			remaining := margin
			if measure > size.Cols {
				remaining = measure.Sub(size.Cols)
			}

			newLines[i] = padding.AddColsPadding(
				remaining, lines[i], opts...
			)
		}

		return newLines, hasNext
	}
}
