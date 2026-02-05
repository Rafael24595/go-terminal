package math_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestAbsSigned(t *testing.T) {
	tests := []struct {
		name string
		in   int64
		out  int64
	}{
		{"positive", 5, 5},
		{"negative", -5, 5},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, math.Abs(tt.in), tt.out)
		})
	}
}

func TestAbsUnsigned(t *testing.T) {
	tests := []struct {
		name string
		in   uint32
	}{
		{"zero", 0},
		{"small", 5},
		{"large", 123456},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, math.Abs(tt.in), tt.in)
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name        string
		val, lo, hi int
		want        int
	}{
		{"inside", 5, 0, 10, 5},
		{"below", -5, 0, 10, 0},
		{"above", 15, 0, 10, 10},
		{"edge_low", 0, 0, 10, 0},
		{"edge_high", 10, 0, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, math.Clamp(tt.val, tt.lo, tt.hi), tt.want)
		})
	}
}

func TestSubClampZeroSigned(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"a>b", 10, 3, 7},
		{"a=b", 5, 5, 0},
		{"a<b", 3, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, math.SubClampZero(tt.a, tt.b), tt.want)
		})
	}
}

func TestSubClampZeroUnsigned(t *testing.T) {
	tests := []struct {
		name string
		a, b uint
		want uint
	}{
		{"a>b", 10, 3, 7},
		{"a=b", 5, 5, 0},
		{"a<b", 3, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, math.SubClampZero(tt.a, tt.b), tt.want)
		})
	}
}

func TestDigitLenSigned(t *testing.T) {
	tests := []struct {
		val  int
		want uint32
	}{
		{0, 1},
		{5, 1},
		{9, 1},
		{10, 2},
		{99, 2},
		{100, 3},
		{-1, 1},
		{-10, 2},
		{-999, 3},
	}

	for _, tt := range tests {
		assert.Equal(t, math.Digits(tt.val), tt.want)
	}
}

func TestDigitsUnsigned(t *testing.T) {
	tests := []struct {
		val  uint64
		want uint32
	}{
		{0, 1},
		{7, 1},
		{42, 2},
		{999, 3},
		{1000, 4},
	}

	for _, tt := range tests {
		assert.Equal(t, math.Digits(tt.val), tt.want)
	}
}
