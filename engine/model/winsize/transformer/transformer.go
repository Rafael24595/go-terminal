package transformer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func WithMargin(rows winsize.Rows, cols uint16) winsize.Transformer {
	return func(w winsize.Winsize) winsize.Winsize {
		return winsize.New(
			math.SubClampZero(w.Rows, rows),
			math.SubClampZero(w.Cols, cols),
		)
	}
}
