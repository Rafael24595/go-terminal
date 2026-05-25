package styler

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

func TestAtomsStyler(t *testing.T) {
	tests := []struct {
		name     string
		atom     style.Atom
		input    string
		expected string
	}{
		{
			name:     "lowercase transformation",
			atom:     style.AtmLower,
			input:    "HelLo GoLang",
			expected: "hello golang",
		},
		{
			name:     "uppercase transformation",
			atom:     style.AtmUpper,
			input:    "HelLo GoLang",
			expected: "HELLO GOLANG",
		},
		{
			name:     "bold does not modify text",
			atom:     style.AtmBold,
			input:    "Hello Golang",
			expected: "Hello Golang",
		},
		{
			name:     "select does not modify text",
			atom:     style.AtmSelect,
			input:    "Hello Golang",
			expected: "Hello Golang",
		},
		{
			name:     "empty string",
			atom:     style.AtmLower,
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := Atoms.Get(tt.atom)
			assert.True(t, ok)

			result := fn(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
