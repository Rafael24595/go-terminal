package winsize

type Transformer func(Winsize) Winsize

// TODO: Use custom type for cols.

type Rows uint16
//type Cols uint16

type Winsize struct {
	Rows Rows
	Cols uint16
}

func NewWinsize(rows Rows, cols uint16) Winsize {
	return Winsize{
		Rows: rows,
		Cols: cols,
	}
}

func (w Winsize) Eq(other Winsize) bool {
	return w.Rows == other.Rows && w.Cols == other.Cols
}
