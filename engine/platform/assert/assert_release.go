//go:build !debug
// +build !debug

package assert

func Unreachable(msg string, a ...any) {}

func True(cond bool, msg string, a ...any) {}

func False(cond bool, msg string, a ...any) {}
