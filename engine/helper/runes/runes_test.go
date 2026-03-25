package runes

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestJoinReverse(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{
			name: "basic",
			in:   []string{"a", "b", "c"},
			out:  "cba",
		},
		{
			name: "words",
			in:   []string{"hello", " ", "golang"},
			out:  "golang hello",
		},
		{
			name: "unicode",
			in:   []string{"🙂", "🚀", "go"},
			out:  "go🚀🙂",
		},
		{
			name: "empty",
			in:   []string{},
			out:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, JoinReverse(tt.in))
		})
	}
}

func TestRuneIndexToByteIndex(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		runeIndex int
		expected  int
		ok        bool
	}{
		{
			name:      "ascii simple",
			text:      "hello",
			runeIndex: 1,
			expected:  1,
			ok:        true,
		},
		{
			name:      "unicode multi-byte",
			text:      "a🙂b",
			runeIndex: 1,
			expected:  1,
			ok:        true,
		},
		{
			name:      "unicode end",
			text:      "a🙂b",
			runeIndex: 3,
			expected:  len("a🙂b"),
			ok:        true,
		},
		{
			name:      "zero index",
			text:      "abc",
			runeIndex: 0,
			expected:  0,
			ok:        true,
		},
		{
			name:      "out of bounds",
			text:      "abc",
			runeIndex: 5,
			expected:  0,
			ok:        false,
		},
		{
			name:      "exact end boundary",
			text:      "abc",
			runeIndex: 3,
			expected:  3,
			ok:        true,
		},
		{
			name:      "empty string",
			text:      "",
			runeIndex: 0,
			expected:  0,
			ok:        true,
		},
		{
			name:      "multi rune unicode",
			text:      "🙂🙂🙂",
			runeIndex: 2,
			expected:  len("🙂🙂"),
			ok:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, ok := RuneIndexToByteIndex(tt.text, tt.runeIndex)

			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.expected, idx)
		})
	}
}

func TestMeasure(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{"ascii", "hello", 5},
		{"unicode", "🙂🙂", 2},
		{"mixed", "a🙂b", 3},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Measure(tt.text))
		})
	}
}

func TestMeasureu(t *testing.T) {
	tests := []struct {
		name string
		text string
		want uint
	}{
		{"ascii", "hello", 5},
		{"unicode", "🙂🙂", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Measureu(tt.text))
		})
	}
}
