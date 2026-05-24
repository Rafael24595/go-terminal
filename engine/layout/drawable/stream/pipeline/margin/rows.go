package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding/options"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func RowsTop(margin winsize.Rows) pipeline.DataTransformer {
	return rowsDataTransformer(margin, options.WithPosition(style.Bottom))
}

func RowsBottom(margin winsize.Rows) pipeline.DataTransformer {
	return rowsDataTransformer(margin, options.WithPosition(style.Top))
}

func RowsMiddle(margin winsize.Rows) pipeline.DataTransformer {
	return rowsDataTransformer(margin, options.WithPosition(style.Middle))
}

func rowsDataTransformer(margin winsize.Rows, opts ...options.RowsOption) pipeline.DataTransformer {
	cfg := options.ResolveRowsConfig(opts...)
	margin = margin * verticalFactor(cfg.Position)

	return func(size winsize.Winsize, _ drawable.Unit, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		linesLen := winsize.Rows(len(lines))

		transformer := padding.Rows(
			hint.Fixed(linesLen+margin),
			opts...,
		)

		return transformer(size, lines), hasNext
	}
}
