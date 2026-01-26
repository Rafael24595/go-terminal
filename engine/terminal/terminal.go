package terminal

type Winsize struct {
	Rows uint16
	Cols uint16
	Err  error
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
