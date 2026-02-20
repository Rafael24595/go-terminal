package wrapper_render

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestPaddingLeft_Strict(t *testing.T) {
	spec := style.SpecPaddingLeft(6, "-")
	cols := 20

	got := paddingLeft(spec, cols, "hi")

	assert.Equal(t, "----hi", got)
	assert.Equal(t, 6, len(got))
}

func TestPaddingLeft_RespectsCols(t *testing.T) {
	spec := style.SpecPaddingLeft(10, "-")
	cols := 5

	got := paddingLeft(spec, cols, "hi")

	assert.Equal(t, "---hi", got)
	assert.Equal(t, 5, len(got))
}

func TestPaddingRight_Strict(t *testing.T) {
	spec := style.SpecPaddingRight(6, ".")
	cols := 20

	got := paddingRight(spec, cols, "hi")

	assert.Equal(t, "hi....", got)
	assert.Equal(t, 6, len(got))
}

func TestPaddingRight_RespectsCols(t *testing.T) {
	spec := style.SpecPaddingRight(10, ".")
	cols := 5

	got := paddingRight(spec, cols, "hi")

	assert.Equal(t, "hi...", got)
	assert.Equal(t, 5, len(got))
}

func TestPaddingCenter_Strict(t *testing.T) {
	spec := style.SpecPaddingCenter(6, "-")
	cols := 20

	got := paddingCenter(spec, cols, "hi")

	assert.Equal(t, "--hi--", got)
	assert.Equal(t, 6, len(got))
}

func TestPaddingCenter_RespectsCols(t *testing.T) {
	spec := style.SpecPaddingCenter(6, "-")
	cols := 4

	got := paddingCenter(spec, cols, "hi")

	assert.Equal(t, "-hi-", got)
	assert.Equal(t, 4, len(got))
}

func TestPaddingCenter_OddSize(t *testing.T) {
	spec := style.SpecPaddingCenter(7, "-")
	cols := 20

	got := paddingCenter(spec, cols, "hi")

	assert.Equal(t, "--hi---", got)
	assert.Equal(t, 7, len(got))
}

func TestRepeatLeft_WithText_Strict(t *testing.T) {
	spec := style.SpecRepeatLeft(3, "-")
	got := repeatLeft(spec, 20, "hi")

	assert.Equal(t, "---hi", got)
}

func TestRepeatLeft_WithoutText_Strict(t *testing.T) {
	spec := style.SpecRepeatLeft(3)
	got := repeatLeft(spec, 20, "ab")

	assert.Equal(t, "bab", got)
}

func TestRepeatRight_WithText_Strict(t *testing.T) {
	spec := style.SpecRepeatRight(3, "-")
	got := repeatRight(spec, 20, "hi")

	assert.Equal(t, "hi---", got)
}

func TestRepeatRight_WithoutText_Strict(t *testing.T) {
	spec := style.SpecRepeatRight(3)
	got := repeatRight(spec, 20, "ab")

	assert.Equal(t, "aba", got)
}

func TestTrimLeft(t *testing.T) {
	tests := []struct {
		name string
		size uint
		in   string
		want string
	}{
		{"trim 2", 2, "golang", "lang"},
		{"trim 1", 1, "zig", "ig"},
		{"trim zero -> min 1", 0, "go", "o"},
		{"empty input", 3, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := style.SpecTrimLeft(tt.size)
			got := trimLeft(spec, tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTrimRight(t *testing.T) {
	tests := []struct {
		name string
		size uint
		in   string
		want string
	}{
		{"trim 2", 2, "golang", "go"},
		{"trim 1", 1, "ziglang", "z"},
		{"trim zero -> min 1", 0, "go", "g"},
		{"empty input", 2, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := style.SpecTrimRight(tt.size)
			got := trimRight(spec, tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFill_Strict(t *testing.T) {
	got := fill(10, 6, "-")

	assert.Equal(t, 6, len(got))
	assert.Equal(t, "------", got)
}

func TestFill_Strict_LongText_Even(t *testing.T) {
	got := fill(20, 10, "go")

	assert.Equal(t, 10, len(got))
	assert.Equal(t, "gogogogogo", got)
}

func TestFill_Strict_LongText_Odd(t *testing.T) {
	got := fill(20, 10, "zig")

	assert.Equal(t, 10, len(got))
	assert.Equal(t, "zigzigzigz", got)
}
