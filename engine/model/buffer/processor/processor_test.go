package processor

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestProcessor_Identity(t *testing.T) {
	input := []rune("Hello Golang")

	buff, facade := Identity(input)

	assert.Equal(t, "Hello Golang", string(buff))
	assert.Equal(t, "Hello Golang", string(facade))
}

func TestProcessor_Number(t *testing.T) {
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
		buff, facade := Number([]rune(c.input))
		assert.Equal(t, c.expected, string(buff))
		assert.Equal(t, c.expected, string(facade))
	}
}

func TestProcessor_Hidden(t *testing.T) {
	input := []rune("password123")

	buff, facade := Hidden(input)

	assert.Equal(t, "password123", string(buff))
	assert.Equal(t, "***********", string(facade))
	assert.Len(t, len(buff), facade)
}

func TestProcessor_Limit(t *testing.T) {
	handler := Limit(6, Identity)

	input := []rune("GolangZiglang")
	buff, facade := handler(input)

	assert.Equal(t, "Golang", string(buff))
	assert.Equal(t, "Golang", string(facade))
	assert.Equal(t, 6, len(buff))
}

func TestProcessor_NumberLimited(t *testing.T) {
	handler := Limit(3, Number)

	input := []rune("1a2b3c4d5")
	buff, _ := handler(input)

	assert.Equal(t, "123", string(buff))
}

func TestProcessor_Inline(t *testing.T) {
	input := []rune("\nHello\nGolang\n")

	buff, facade := Inline(input)

	assert.Equal(t, " Hello Golang ", string(buff))
	assert.Equal(t, " Hello Golang ", string(facade))
}
