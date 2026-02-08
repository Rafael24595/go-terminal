//go:build debug
// +build debug

package assert

import "fmt"

func Unreachable(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}

func AssertTrue(cond bool, format string, a ...any) {
	if !cond {
		msg := fmt.Sprintf(format, a...)
		panic(msg)
	}
}

func AssertFalse(cond bool, format string, a ...any) {
	if cond {
		msg := fmt.Sprintf(format, a...)
		panic(msg)
	}
}
