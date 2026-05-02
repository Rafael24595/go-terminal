package transformer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func WithMargin(rows winsize.Rows, cols winsize.Cols) winsize.Transformer {
	return func(w winsize.Winsize) winsize.Winsize {
		return winsize.New(
			w.Rows.Clamp(rows),
			w.Cols.Clamp(cols),
		)
	}
}
