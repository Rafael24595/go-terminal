package math

import (
	"cmp"
)

type Number interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Abs[T Number](val T) T {
	switch v := any(val).(type) {
	case int:
		if v < 0 {
			return T(-v)
		}
	case int8:
		if v < 0 {
			return T(-v)
		}
	case int16:
		if v < 0 {
			return T(-v)
		}
	case int32:
		if v < 0 {
			return T(-v)
		}
	case int64:
		if v < 0 {
			return T(-v)
		}
	}
	return val
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

func Digits[T Number](val T) uint32 {
	if val == 0 {
		return 1
	}

	fix := Abs(val) / 10
	count := uint32(1)

	for fix > 0 {
		fix = fix / 10
		count++
	}

	return count
}
