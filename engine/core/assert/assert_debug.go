//go:build debug
// +build debug

package assert

import "fmt"

func AssertTrue(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func AssertFalse(cond bool, msg string) {
	if cond {
		panic(msg)
	}
}

func AssertfTrue(cond bool, format string, a ...any) {
	if !cond {
		msg := fmt.Sprintf(format, a...)
		panic(msg)
	}
}

func AssertfFalse(cond bool, format string, a ...any) {
	if cond {
		msg := fmt.Sprintf(format, a...)
		panic(msg)
	}
}
