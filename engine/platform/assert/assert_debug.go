//go:build debug
// +build debug

package assert

import "fmt"

func Unreachable(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}

func True(cond bool, format string, a ...any) {
	if !cond {
		msg := fmt.Sprintf(format, a...)
		panic(msg)
	}
}

func False(cond bool, format string, a ...any) {
	if cond {
		msg := fmt.Sprintf(format, a...)
		panic(msg)
	}
}
