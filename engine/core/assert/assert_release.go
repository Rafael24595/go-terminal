//go:build !debug
// +build !debug

package assert

func AssertTrue(cond bool, msg string) {}

func AssertFalse(cond bool, msg string) {}

func AssertfFalse(cond bool, format string, a ...any) {}

func AssertfTrue(cond bool, format string, a ...any) {}
