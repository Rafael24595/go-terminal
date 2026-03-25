package delta

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestDeltaMeasure(t *testing.T) {
	d := Delta{
		Text: "a🙂b",
	}

	assert.Equal(t, uint(3), d.Measure())
}

func TestApplyBasic(t *testing.T) {
	buffer := []rune("hello")

	d := &Delta{
		Start: 0,
		End:   0,
		Text:  "X",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "Xhello", string(result))
}

func TestApplyReplaceMiddle(t *testing.T) {
	buffer := []rune("hello")

	d := &Delta{
		Start: 1,
		End:   4,
		Text:  "i",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "hio", string(result))
}

func TestApplyDelete(t *testing.T) {
	buffer := []rune("hello")

	d := &Delta{
		Start: 1,
		End:   4,
		Text:  "",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "ho", string(result))
}

func TestApplyInsertMiddle(t *testing.T) {
	buffer := []rune("helo")

	d := &Delta{
		Start: 2,
		End:   2,
		Text:  "l",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "hello", string(result))
}

func TestApplyInsertEnd(t *testing.T) {
	buffer := []rune("hello")

	d := &Delta{
		Start: 5,
		End:   5,
		Text:  "!",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "hello!", string(result))
}

func TestApplyOutOfBounds(t *testing.T) {
	buffer := []rune("hello")

	d := &Delta{
		Start: 10,
		End:   12,
		Text:  "X",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "hello", string(result))
}

func TestApplyUnicode(t *testing.T) {
	buffer := []rune("a🙂b")

	d := &Delta{
		Start: 1,
		End:   2,
		Text:  "🚀",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "a🚀b", string(result))
}

func TestApplyReplaceAll(t *testing.T) {
	buffer := []rune("hello")

	d := &Delta{
		Start: 0,
		End:   5,
		Text:  "bye",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "bye", string(result))
}

func TestApplyEmptyBuffer(t *testing.T) {
	buffer := []rune("")

	d := &Delta{
		Start: 0,
		End:   0,
		Text:  "hi",
	}

	result := Apply(buffer, d)

	assert.Equal(t, "hi", string(result))
}
