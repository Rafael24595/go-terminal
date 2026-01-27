//go:build !mock_cmd && windows
// +build !mock_cmd,windows

package wrapper_terminal

import (
	"syscall"
	"unsafe"
)

const (
	INPUT_KEYBOARD  = 1
	KEYEVENTF_KEYUP = 0x0002
)

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	sendInputProc = user32.NewProc("SendInput")
)

type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uint64
}

type input struct {
	inputType uint32
	ki        keyboardInput
	padding   uint64
}

func sendDummyKey() {
	inputs := []input{
		{
			inputType: INPUT_KEYBOARD,
			ki: keyboardInput{
				wVk: VK_SHIFT,
			},
		},
		{
			inputType: INPUT_KEYBOARD,
			ki: keyboardInput{
				wVk:     VK_SHIFT,
				dwFlags: KEYEVENTF_KEYUP,
			},
		},
	}

	ret, _, err := sendInputProc.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(inputs[0])),
	)

	if ret == 0 {
		//TODO: Log and force user manual input instead panic.
		panic("SendInput failed:" + err.Error())
	}
}
