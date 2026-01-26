//go:build !mock_cmd && windows
// +build !mock_cmd,windows

package wrapper_console

import (
	"syscall"
	"unsafe"

	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type coord struct {
	X int16
	Y int16
}

type smallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type consoleScreenBufferInfo struct {
	Size              coord
	CursorPosition    coord
	Attributes        uint16
	Window            smallRect
	MaximumWindowSize coord
}

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	getCSBI  = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

func Size() terminal.Winsize {
	handle := syscall.Handle(syscall.Stdout)

	var info consoleScreenBufferInfo
	r, _, e := getCSBI.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&info)),
	)

	if r == 0 {
		return terminal.Winsize{
			Err: e,
		}
	}

	return terminal.Winsize{
		Rows: uint16(info.Window.Bottom-info.Window.Top) + 1,
		Cols: uint16(info.Window.Right-info.Window.Left) + 1,
	}
}
