package terminal

import (
	"github.com/Rafael24595/go-terminal/engine/model/key"
)

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
