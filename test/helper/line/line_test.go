package line_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/helper/line"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestFindLineStart_Simple(t *testing.T) {
	buf := []rune("abc\ndef")

	assert.Equal(t, 0, line.FindLineStart(buf, 0))
	assert.Equal(t, 0, line.FindLineStart(buf, 2))
	assert.Equal(t, 4, line.FindLineStart(buf, 4))
	assert.Equal(t, 4, line.FindLineStart(buf, 6))
}

func TestFindNextLineStart_Simple(t *testing.T) {
	buf := []rune("abc\ndef")

	assert.Equal(t, 4, line.FindNextLineStart(buf, 0))
	assert.Equal(t, 4, line.FindNextLineStart(buf, 2))
	assert.Equal(t, -1, line.FindNextLineStart(buf, 4))
}

func TestFindPrevLineStart_EmptyLine(t *testing.T) {
	buf := []rune("abc\n\ndef")

	from := 5

	expected := 4

	got := line.FindPrevLineStart(buf, from)

	assert.Equal(t, expected, got)
}

func TestFindPrevLineStart_MultipleEmptyLines(t *testing.T) {
	buf := []rune("a\n\n\nb")

	from := 4

	expected := 3

	got := line.FindPrevLineStart(buf, from)

	assert.Equal(t, expected, got)
}

func TestFindPrevLineStart_Normal(t *testing.T) {
	buf := []rune("a\nb\nc")

	from := 4

	expected := 2
	got := line.FindPrevLineStart(buf, from)

	assert.Equal(t, expected, got)
}

func TestClampToLine_EmptyLine(t *testing.T) {
	buf := []rune("abc\n\ndef")

	lineStart := 4
	col := 10

	assert.Equal(t, 4, line.ClampToLine(buf, lineStart, col))
}
