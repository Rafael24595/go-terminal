//go:build !mock_cmd && windows
// +build !mock_cmd,windows

package wrapper_terminal

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

const (
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = uint32(0x0004)
	ENABLE_PROCESSED_INPUT             = uint32(0x0004)
	ENABLE_LINE_INPUT                  = uint32(0x0002)
	ENABLE_ECHO_INPUT                  = uint32(0x0001)
	ENABLE_VIRTUAL_TERMINAL_INPUT      = uint32(0x0200)
)

func onStart() (uintptr, error) {
	sendDummyKey()
	enableANSI()
	return enableRaw()
}

func onClose(rawmode uintptr) {
	restoreRaw(rawmode)
}

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

func enableANSI() {
	handle := syscall.Handle(syscall.Stdout)

	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	var mode uint32
	getConsoleMode.Call(uintptr(handle),
		uintptr(unsafe.Pointer(&mode)))

	setConsoleMode.Call(uintptr(handle),
		uintptr(mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING))
}

func enableRaw() (uintptr, error) {
	handle := syscall.Handle(syscall.Stdin)

	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	var mode uint32
	getConsoleMode.Call(uintptr(handle),
		uintptr(unsafe.Pointer(&mode)))

	oldMode := mode
	mode &^= ENABLE_PROCESSED_INPUT |
		ENABLE_LINE_INPUT |
		ENABLE_ECHO_INPUT

	mode |= ENABLE_VIRTUAL_TERMINAL_INPUT

	setConsoleMode.Call(uintptr(handle), uintptr(mode))

	return uintptr(oldMode), nil
}

func restoreRaw(old uintptr) {
	handle := syscall.Handle(syscall.Stdin)
	kernel32.NewProc("SetConsoleMode").Call(uintptr(handle), old)
}
