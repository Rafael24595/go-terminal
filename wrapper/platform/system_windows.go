//go:build !mock_cmd && windows
// +build !mock_cmd,windows

package platform

import (
	"context"
	"syscall"
	"time"
	"unsafe"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
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

func OnStart() (uintptr, error) {
	err := sendDummyKey()
	if err != nil {
		return 0, err
	}

	err = enableANSI()
	if err != nil {
		return 0, err
	}

	return enableRaw()
}

func OnClose(rawmode uintptr) error {
	return restoreRaw(rawmode)
}

func Size() (winsize.Winsize, error) {
	handle := syscall.Handle(syscall.Stdout)

	var info consoleScreenBufferInfo
	ret, _, err := getCSBI.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&info)),
	)

	if ret == 0 {
		return winsize.Winsize{}, err
	}

	return winsize.New(
		winsize.Rows(info.Window.Bottom-info.Window.Top)+1,
		uint16(info.Window.Right-info.Window.Left)+1,
	), nil
}

func ResizeSystemEvents(ctx context.Context, drt time.Duration) <-chan winsize.Winsize {
	return ResizeReactiveEvents(ctx, drt)
}

func enableANSI() error {
	handle := syscall.Handle(syscall.Stdout)

	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	var mode uint32
	ret, _, err := getConsoleMode.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&mode)),
	)

	if ret == 0 {
		return err
	}

	ret, _, err = setConsoleMode.Call(
		uintptr(handle),
		uintptr(mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING),
	)

	if ret == 0 {
		return err
	}

	return nil
}

func enableRaw() (uintptr, error) {
	handle := syscall.Handle(syscall.Stdin)

	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	var mode uint32
	ret, _, err := getConsoleMode.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&mode)),
	)

	if ret == 0 {
		return 0, err
	}

	oldMode := mode
	mode &^= ENABLE_PROCESSED_INPUT |
		ENABLE_LINE_INPUT |
		ENABLE_ECHO_INPUT

	mode |= ENABLE_VIRTUAL_TERMINAL_INPUT

	ret, _, err = setConsoleMode.Call(
		uintptr(handle),
		uintptr(mode),
	)

	if ret == 0 {
		return 0, err
	}

	return uintptr(oldMode), nil
}

func restoreRaw(old uintptr) error {
	handle := syscall.Handle(syscall.Stdin)

	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	ret, _, err := setConsoleMode.Call(
		uintptr(handle),
		old,
	)

	if ret == 0 {
		return err
	}

	return nil
}
