//go:build !mock_cmd && linux
// +build !mock_cmd,linux

package wrapper_console

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func Size() terminal.Winsize {
	ws := &winsize{}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if errno != 0 {
		return terminal.Winsize{
			Err: errno,
		}
	}

	return terminal.Winsize{
		Rows: uint16(ws.Row),
		Cols: uint16(ws.Col),
	}
}
