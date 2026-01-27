//go:build !mock_cmd && linux
// +build !mock_cmd,linux

package wrapper_terminal

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const (
	ICANON = 0x0001
	ECHO   = 0x0002
)

func onStart() (uintptr, error) {
	return enableRaw()
}

func onClose(rawmode uintptr) {
	restoreRaw(rawmode)
}

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

type termios struct {
	Iflag, Oflag, Cflag, Lflag uint32
	Cc                         [20]byte
	Ispeed, Ospeed             uint32
}

func enableRaw() (uintptr, error) {
	fd := uintptr(syscall.Stdin)

	var tio termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&tio)), 0, 0, 0)
	if err != 0 {
		return 0, err
	}

	tio.Lflag &^= ICANON | ECHO

	tio.Cc[syscall.VMIN] = 0
	tio.Cc[syscall.VTIME] = 1

	_, _, err = syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&tio)), 0, 0, 0)
	if err != 0 {
		return 0, err
	}

	return uintptr(unsafe.Pointer(&tio)), nil
}

func restoreRaw(old uintptr) {
	fd := os.Stdin.Fd()
	syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), old, 0, 0, 0)
}
