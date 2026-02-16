package helper_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestNumberToAlpha(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"invalid zero", 0, "?"},
		{"invalid negative", -5, "?"},
		{"single letter a", 1, "a"},
		{"single letter b", 2, "b"},
		{"single letter z", 26, "z"},
		{"double letter aa", 27, "aa"},
		{"double letter ab", 28, "ab"},
		{"double letter az", 52, "az"},
		{"double letter ba", 53, "ba"},
		{"triple letter aaa", 703, "aaa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.NumberToAlpha(tt.input)
			assert.Equal(t, result, tt.expected)
		})
	}
}
