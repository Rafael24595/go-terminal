package math

import (
	"cmp"
)

type Number interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

func Clamp[T cmp.Ordered](val, lower, upper T) T {
	return max(lower, min(val, upper))
}

func SubClampZero[T Number](a, b T) T {
	if a < b {
		return 0
	}
	return a - b
}
