package wrapper_render

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

func TestPaddingLeft_Strict(t *testing.T) {
	spec := style.SpecPaddingLeft(6, "-")
	cols := 20

	text := "hi"
	size := runes.Measure(text)

	got := paddingLeft(spec, cols, text, size)

	assert.Equal(t, "----hi", got)
	assert.Equal(t, 6, len(got))
}

func TestPaddingLeft_RespectsCols(t *testing.T) {
	spec := style.SpecPaddingLeft(10, "-")
	cols := 5

	text := "hi"
	size := runes.Measure(text)

	got := paddingLeft(spec, cols, text, size)

	assert.Equal(t, "---hi", got)
	assert.Equal(t, 5, len(got))
}

func TestPaddingRight_Strict(t *testing.T) {
	spec := style.SpecPaddingRight(6, ".")
	cols := 20

	text := "hi"
	size := runes.Measure(text)

	got := paddingRight(spec, cols, text, size)

	assert.Equal(t, "hi....", got)
	assert.Equal(t, 6, len(got))
}

func TestPaddingRight_RespectsCols(t *testing.T) {
	spec := style.SpecPaddingRight(10, ".")
	cols := 5

	text := "hi"
	size := runes.Measure(text)

	got := paddingRight(spec, cols, text, size)

	assert.Equal(t, "hi...", got)
	assert.Equal(t, 5, len(got))
}

func TestPaddingCenter_Strict(t *testing.T) {
	spec := style.SpecPaddingCenter(6, "-")
	cols := 20

	text := "hi"
	size := runes.Measure(text)

	got := paddingCenter(spec, cols, text, size)

	assert.Equal(t, "--hi--", got)
	assert.Equal(t, 6, len(got))
}

func TestPaddingCenter_RespectsCols(t *testing.T) {
	spec := style.SpecPaddingCenter(6, "-")
	cols := 4

	text := "hi"
	size := runes.Measure(text)

	got := paddingCenter(spec, cols, text, size)

	assert.Equal(t, "-hi-", got)
	assert.Equal(t, 4, len(got))
}

func TestPaddingCenter_OddSize(t *testing.T) {
	spec := style.SpecPaddingCenter(7, "-")
	cols := 20

	text := "hi"
	size := runes.Measure(text)

	got := paddingCenter(spec, cols, text, size)

	assert.Equal(t, "--hi---", got)
	assert.Equal(t, 7, len(got))
}

func TestRepeatLeft_WithText_Strict(t *testing.T) {
	spec := style.SpecRepeatLeft(3, "-")

	text := "hi"
	size := runes.Measure(text)

	got := repeatLeft(spec, 20, text, size)

	assert.Equal(t, "---hi", got)
}

func TestRepeatLeft_WithoutText_Strict(t *testing.T) {
	spec := style.SpecRepeatLeft(3)

	text := "ab"
	size := runes.Measure(text)

	got := repeatLeft(spec, 20, text, size)

	assert.Equal(t, "bab", got)
}

func TestRepeatRight_WithText_Strict(t *testing.T) {
	spec := style.SpecRepeatRight(3, "-")

	text := "hi"
	size := runes.Measure(text)

	got := repeatRight(spec, 20, text, size)

	assert.Equal(t, "hi---", got)
}

func TestRepeatRight_WithoutText_Strict(t *testing.T) {
	spec := style.SpecRepeatRight(3)

	text := "ab"
	size := runes.Measure(text)

	got := repeatRight(spec, 20, text, size)

	assert.Equal(t, "aba", got)
}

func TestTrimLeft_Standard(t *testing.T) {
	tests := []struct {
		name string
		size uint
		in   string
		want string
	}{
		{
			name: "keep last 2 characters",
			size: 2,
			in:   "golang",
			want: "ng",
		},
		{
			name: "keep last character",
			size: 1,
			in:   "zig",
			want: "g",
		},
		{
			name: "fallback to minimum 1 when size is 0",
			size: 0,
			in:   "go",
			want: "o",
		},
		{
			name: "handle empty string",
			size: 3,
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := style.SpecTrimLeft(tt.size)

			size := runes.Measure(tt.in)
			got := trimLeft(spec, tt.in, size)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && size > 0 {
				assert.Equal(t, tt.size, runes.Measureu(got))
			}
		})
	}
}

func TestTrimLeft_WithEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		size     uint
		ellipsis string
		in       string
		want     string
	}{
		{
			name:     "prepend ellipsis when space allows",
			size:     5,
			ellipsis: ".",
			in:       "golang",
			want:     "...ng",
		},
		{
			name:     "skip ellipsis if it consumes too much space",
			size:     2,
			ellipsis: ".",
			in:       "ziglang",
			want:     "ng",
		},
		{
			name:     "bypass ellipsis logic when size+elipSize exceeds logical limits",
			size:     1,
			ellipsis: "..",
			in:       "rust",
			want:     "t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := style.SpecTrimTextLeft(tt.size, tt.ellipsis)
			size := runes.Measure(tt.in)

			got := trimLeft(spec, tt.in, size)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && size > 0 {
				assert.Equal(t, tt.size, runes.Measureu(got))
			}
		})
	}
}

func TestTrimRight_Standard(t *testing.T) {
	tests := []struct {
		name string
		size uint
		in   string
		want string
	}{
		{
			name: "keep first 2 characters",
			size: 2,
			in:   "golang",
			want: "go",
		},
		{
			name: "keep first character",
			size: 1,
			in:   "ziglang",
			want: "z"},
		{
			name: "fallback to minimum 1 when size is 0",
			size: 0,
			in:   "go",
			want: "g",
		},
		{
			name: "handle empty string",
			size: 2,
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := style.SpecTrimRight(tt.size)

			size := runes.Measure(tt.in)
			got := trimRight(spec, tt.in, size)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && size > 0 {
				assert.Equal(t, tt.size, runes.Measureu(got))
			}
		})
	}
}

func TestTrimRight_WithEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		size     uint
		ellipsis string
		in       string
		want     string
	}{
		{
			name:     "append ellipsis when space allows",
			size:     5,
			ellipsis: ".",
			in:       "golang",
			want:     "go...",
		},
		{
			name:     "skip ellipsis and return raw trim when space is tight",
			size:     2,
			ellipsis: ".",
			in:       "ziglang",
			want:     "zi",
		},
		{
			name:     "return direct trim when logical size is exceeded",
			size:     1,
			ellipsis: "...",
			in:       "test",
			want:     "t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := style.SpecTrimTextRight(tt.size, tt.ellipsis)
			size := runes.Measure(tt.in)

			got := trimRight(spec, tt.in, size)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && size > 0 {
				assert.Equal(t, tt.size, runes.Measureu(got))
			}
		})
	}
}

func TestFill_Strict(t *testing.T) {
	text := "-"
	size := runes.Measure(text)

	spec := style.SpecFill(10)
	got := fill(spec, 6, text, size)

	assert.Equal(t, 6, len(got))
	assert.Equal(t, "------", got)
}

func TestFill_Strict_LongText_Even(t *testing.T) {
	text := "go"
	size := runes.Measure(text)

	spec := style.SpecFill(20)
	got := fill(spec, 10, text, size)

	assert.Equal(t, 10, len(got))
	assert.Equal(t, "gogogogogo", got)
}

func TestFill_Strict_LongText_Odd(t *testing.T) {
	text := "zig"
	size := runes.Measure(text)

	spec := style.SpecFill(20)
	got := fill(spec, 10, text, size)

	assert.Equal(t, 10, len(got))
	assert.Equal(t, "zigzigzigz", got)
}
