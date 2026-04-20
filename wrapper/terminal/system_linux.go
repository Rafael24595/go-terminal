//go:build !mock_cmd && linux
// +build !mock_cmd,linux

package wrapper_terminal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/Rafael24595/go-terminal/engine/model/winsize"
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

type linuxWinsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func Size() (winsize.Winsize, error) {
	ws := &linuxWinsize{}

	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if err != 0 {
		return winsize.Winsize{}, err
	}

	return winsize.New(
		winsize.Rows(ws.Row),
		uint16(ws.Col),
	), nil
}

func ResizeEvents(ctx context.Context) chan winsize.Winsize {
	out := make(chan winsize.Winsize, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)

	go func() {
		defer close(out)
		defer signal.Stop(sig)

		for {
			select {
			case <-ctx.Done():
				return
			case <-sig:
				size, err := Size()
				if err != nil {
					continue
				}

				select {
				case out <- size:
				default:
				}
			}
		}
	}()

	return out
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
