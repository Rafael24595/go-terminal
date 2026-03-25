package buffer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestHandler_String(t *testing.T) {
	handler := NewRuneHandler(String)
	input := []rune("Hello Golang")

	buff, facade := handler(input)

	assert.Equal(t, "Hello Golang", string(buff))
	assert.Equal(t, "Hello Golang", string(facade))
}

func TestHandler_Number(t *testing.T) {
	handler := NewRuneHandler(Number)

	cases := []struct {
		input    string
		expected string
	}{
		{"123a45", "12345"},
		{"123456", "123456"},
		{"-12.3", "-12.3"},
		{"abc!@#", ""},
		{"12,45", "12,45"},
		{"-12.-3", "-12.3"},
		{"12.3.4", "12.34"},
	}

	for _, c := range cases {
		buff, facade := handler([]rune(c.input))
		assert.Equal(t, c.expected, string(buff))
		assert.Equal(t, c.expected, string(facade))
	}
}

func TestHandler_Hidden(t *testing.T) {
	handler := NewRuneHandler(Hidden)
	input := []rune("password123")

	buff, facade := handler(input)

	assert.Equal(t, "password123", string(buff))
	assert.Equal(t, "***********", string(facade))
	assert.Equal(t, len(buff), len(facade))
}

func TestHandler_Limit(t *testing.T) {
	limit := uint64(6)
	handler := NewLimitedRuneHandler(limit, String)

	input := []rune("GolangZiglang")
	buff, facade := handler(input)

	assert.Equal(t, "Golang", string(buff))
	assert.Equal(t, "Golang", string(facade))
	assert.Equal(t, 6, len(buff))
}

func TestHandler_NumberLimited(t *testing.T) {
	handler := NewLimitedRuneHandler(3, Number)
	
	input := []rune("1a2b3c4d5")
	buff, _ := handler(input)

	assert.Equal(t, "123", string(buff))
}
