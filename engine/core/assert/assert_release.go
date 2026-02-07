//go:build !debug
// +build !debug

package assert

func Unreachable(msg string, a ...any) {}

func AssertTrue(cond bool, msg string, a ...any) {}

func AssertFalse(cond bool, msg string, a ...any) {}
