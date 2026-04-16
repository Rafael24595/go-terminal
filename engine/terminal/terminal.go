package terminal

import (
	"github.com/Rafael24595/go-terminal/engine/model/key"
)

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

type Terminal struct {
	OnStart   func() error
	OnClose   func() error
	Size      func() (Winsize, error)
	Clear     func() error
	ReadKey   func() (*key.Key, error)
	Write     func(string) error
	WriteLine func(...string) error
	WriteAll  func(string) error
	Flush     func() error
}
