package terminal

type Winsize struct {
	Rows uint16
	Cols uint16
	Err  error
}

func NewWinsize(rows, cols uint16) Winsize {
	return Winsize{
		Rows: rows,
		Cols: cols,
	}
}

type Terminal struct {
	OnStart   func() error
	OnClose   func() error
	Size      func() Winsize
	Clear     func() error
	Write     func(string) error
	WriteLine func(...string) error
	WriteAll  func(string) error
	Flush     func() error
}
