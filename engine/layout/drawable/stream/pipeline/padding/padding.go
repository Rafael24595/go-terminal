package padding

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

func Rows(rows hint.Size[winsize.Rows], opts ...options.RowsOption) pipeline.DataTransformer {
	transformer := padding.Rows(rows, opts...)
	return func(size winsize.Winsize, _ drawable.Unit, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		return transformer(size, lines), hasNext
	}
}

func Cols(cols hint.Size[winsize.Cols], position ...style.HorizontalPosition) pipeline.DataTransformer {
	return ColsWithDefault(cols, padding.DefaultColFrag, position...)
}

func ColsWithDefault(cols hint.Size[winsize.Cols], frag string, position ...style.HorizontalPosition) pipeline.DataTransformer {
	transformer := padding.ColsWithDefault(cols, frag, position...)
	return func(size winsize.Winsize, _ drawable.Unit, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		return transformer(size, lines), hasNext
	}
}
