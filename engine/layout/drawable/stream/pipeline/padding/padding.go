package padding

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func RowsDataTransformer(rows padding.SizeHint[winsize.Rows], position ...style.VerticalPosition) pipeline.DataTransformer {
	transformer := padding.Rows(rows, position...)
	return func(size winsize.Winsize, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		return transformer(size, lines), hasNext
	}
}

func ColsDrawTransformer(cols padding.SizeHint[winsize.Cols], position ...style.HorizontalPosition) pipeline.DataTransformer {
	transformer := padding.Cols(cols, position...)
	return func(size winsize.Winsize, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		return transformer(size, lines), hasNext
	}
}
