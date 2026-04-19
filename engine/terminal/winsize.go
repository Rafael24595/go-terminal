package terminal

// TODO: Use custom type for rows and cols.

//type Rows uint16
//type Cols uint16

type Winsize struct {
	Rows uint16
	Cols uint16
}

func NewWinsize(rows, cols uint16) Winsize {
	return Winsize{
		Rows: rows,
		Cols: cols,
	}
}

func (w Winsize) Eq(other Winsize) bool {
	return w.Rows == other.Rows && w.Cols == other.Cols
}
