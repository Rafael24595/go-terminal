package winsize

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

type Transformer func(Winsize) Winsize

type Rows uint16

func (r Rows) Clamp(o Rows) Rows {
	return math.SubClampZero(r, o)
}

type Cols uint16

func (c Cols) Clamp(o Cols) Cols {
	return math.SubClampZero(c, o)
}

type Winsize struct {
	Rows Rows
	Cols Cols
}

func New(rows Rows, cols Cols) Winsize {
	return Winsize{
		Rows: rows,
		Cols: cols,
	}
}

func (w Winsize) Eq(other Winsize) bool {
	return w.Rows == other.Rows && w.Cols == other.Cols
}
