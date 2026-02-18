package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestArgumentNumericConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want int64
	}{
		{"int to int64", int(1), 1},
		{"int8 to int64", int8(1), 1},
		{"int16 to int64", int16(1), 1},
		{"int32 to int64", int32(1), 1},
		{"int64 to int64", int64(1), 1},

		{"uint to int64", uint(1), 1},
		{"uint8 to int64", uint8(1), 1},
		{"uint16 to int64", uint16(1), 1},
		{"uint32 to int64", uint32(1), 1},
		{"uint64 to int64", uint64(1), 1},

		{"float32 to int64", float32(1.0), 1},
		{"float64 to int64", float64(1.0), 1},

		{"bool true to int64", true, 1},
		{"bool false to int64", false, 0},

		{"string '123' to int64", "123", 123},
		{"string '0' to int64", "0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ArgumentFrom(tt.from).Int64()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentUnsignedConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want uint64
	}{
		{"int to uint64", 42, 42},
		{"uint to uint64", uint(42), 42},
		{"uint8 to uint64", uint8(42), 42},
		{"uint16 to uint64", uint16(42), 42},
		{"uint32 to uint64", uint32(42), 42},
		{"uint64 to uint64", uint64(42), 42},

		{"bool true to uint64", true, 1},
		{"bool false to uint64", false, 0},

		{"string '42' to uint64", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ArgumentFrom(tt.from).Uint64()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentFloatConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want float64
	}{
		{"int to float64", 42, 42},
		{"int64 to float64", int64(42), 42},

		{"float32 to float64", float32(42.0), 42},
		{"float64 to float64", float64(42.0), 42},

		{"bool true to float64", true, 1},
		{"bool false to float64", false, 0},

		{"string '42.5' to float64", "42.5", 42.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ArgumentFrom(tt.from).Float64()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentBoolConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want bool
	}{
		{"bool true", true, true},
		{"bool false", false, false},

		{"string true", "true", true},
		{"string false", "false", false},

		{"int zero", 0, false},
		{"int nonzero", 42, true},

		{"float zero", float64(0), false},
		{"float nonzero", float64(3.14), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ArgumentFrom(tt.from).Bool()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentStringConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want string
	}{
		{"nil", nil, ""},
		
		{"string", "hello", "hello"},

		{"int", 42, "42"},

		{"float32", float32(3.14), "3.14"},
		{"float64", float64(3.14), "3.14"},

		{"bool true", true, "true"},
		{"bool false", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArgumentFrom(tt.from).String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentDefaults(t *testing.T) {
	tests := []struct {
		name string
		from any
		def  any
		want any
		fn   func(*Argument) any
	}{
		{"Intd valid int", 42, 100, 42, func(a *Argument) any { return a.Intd(100) }},
		{"Intd invalid string", "abc", 100, 100, func(a *Argument) any { return a.Intd(100) }},

		{"Int64d valid int64", int64(99), 100, int64(99), func(a *Argument) any { return a.Int64d(100) }},
		{"Int64d invalid string", "xyz", 100, int64(100), func(a *Argument) any { return a.Int64d(100) }},

		{"Uint64d valid uint64", uint64(77), 100, uint64(77), func(a *Argument) any { return a.Uint64d(100) }},
		{"Uint64d invalid string", "xyz", 100, uint64(100), func(a *Argument) any { return a.Uint64d(100) }},

		{"Float32d valid float", float32(3.14), 2.71, float32(3.14), func(a *Argument) any { return a.Float32d(2.71) }},
		{"Float32d invalid string", "abc", 2.71, float32(2.71), func(a *Argument) any { return a.Float32d(2.71) }},

		{"Float64d valid float", float64(1.618), 3.14, float64(1.618), func(a *Argument) any { return a.Float64d(3.14) }},
		{"Float64d invalid string", "abc", 3.14, float64(3.14), func(a *Argument) any { return a.Float64d(3.14) }},

		{"Boold valid true", true, false, true, func(a *Argument) any { return a.Boold(false) }},
		{"Boold valid false", false, true, false, func(a *Argument) any { return a.Boold(true) }},
		{"Boold invalid string", "notbool", true, true, func(a *Argument) any { return a.Boold(true) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := ArgumentFrom(tt.from)
			got := tt.fn(arg)
			assert.Equal(t, tt.want, got)
		})
	}
}
