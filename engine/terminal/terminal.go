package terminal

import (
	"github.com/Rafael24595/go-terminal/engine/model/key"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
)

type Terminal struct {
	OnStart      func() error
	OnClose      func() error
	ResizeEvents func() <-chan winsize.Winsize
	KeyEvents    func() <-chan key.Key
	Size         func() (winsize.Winsize, error)
	Clear        func() error
	Write        func(...string) error
	WriteAll     func(string) error
	Flush        func() error
}
