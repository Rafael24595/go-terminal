package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func RowsTop(margin winsize.Rows) pipeline.DataTransformer {
	return rowsDataTransformer(margin, *padding.DefaultRowFrag(), style.Bottom)
}

func RowsBottom(margin winsize.Rows) pipeline.DataTransformer {
	return rowsDataTransformer(margin, *padding.DefaultRowFrag(), style.Top)
}

func RowsMiddle(margin winsize.Rows) pipeline.DataTransformer {
	return rowsDataTransformer(margin, *padding.DefaultRowFrag(), style.Middle)
}

func rowsDataTransformer(margin winsize.Rows, frag text.Fragment, position style.VerticalPosition) pipeline.DataTransformer {
	margin = margin * verticalFactor(position)

	return func(size winsize.Winsize, _ drawable.Unit, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		linesLen := winsize.Rows(len(lines))

		transformer := padding.RowsWithDefault(
			hint.Fixed(linesLen+margin),
			frag,
			position,
		)

		return transformer(size, lines), hasNext
	}
}
