package key_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/key"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestModMask_HasAny(t *testing.T) {
	tests := []struct {
		name     string
		mask     key.ModMask
		mods     []key.ModMask
		expected bool
	}{
		{
			name:     "no mods present",
			mask:     0,
			mods:     []key.ModMask{key.ModAlt},
			expected: false,
		},
		{
			name:     "single mod present",
			mask:     key.ModAlt,
			mods:     []key.ModMask{key.ModAlt},
			expected: true,
		},
		{
			name:     "one of many mods present",
			mask:     key.ModAlt,
			mods:     []key.ModMask{key.ModAlt, key.ModCtrl},
			expected: true,
		},
		{
			name:     "none of the mods present",
			mask:     key.ModShift,
			mods:     []key.ModMask{key.ModAlt, key.ModCtrl},
			expected: false,
		},
		{
			name:     "multiple mods present in mask",
			mask:     key.ModAlt | key.ModCtrl,
			mods:     []key.ModMask{key.ModAlt, key.ModCtrl},
			expected: true,
		},
		{
			name:     "no arguments",
			mask:     key.ModAlt,
			mods:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mask.HasAny(tt.mods...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestModMask_HasNone(t *testing.T) {
	tests := []struct {
		name     string
		mask     key.ModMask
		mods     []key.ModMask
		expected bool
	}{
		{
			name:     "none present",
			mask:     key.ModAlt,
			mods:     []key.ModMask{key.ModShift, key.ModCtrl},
			expected: true,
		},
		{
			name:     "one present",
			mask:     key.ModAlt,
			mods:     []key.ModMask{key.ModAlt},
			expected: false,
		},
		{
			name:     "one of many present",
			mask:     key.ModAlt | key.ModShift,
			mods:     []key.ModMask{key.ModCtrl, key.ModAlt},
			expected: false,
		},
		{
			name:     "no arguments",
			mask:     key.ModAlt,
			mods:     nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mask.HasNone(tt.mods...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

