package line_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-terminal/engine/helper/line"
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

func TestFindLineStart(t *testing.T) {
	buf := []rune("line1\nline2\nline3")

	tests := []struct {
		name string
		from int
		want int
	}{
		{"middle of line2", 8, 6},
		{"start of line3", 12, 12},
		{"start of buffer", 0, 0},
		{"middle of line1", 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := line.FindLineStart(buf, tt.from)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDistanceFromLF(t *testing.T) {
	buf := []rune("line1\nline2\nline3")

	tests := []struct {
		name string
		from int
		want int
	}{
		{"middle of line1", 3, 3},
		{"end of line1", 5, 5},
		{"middle of line2", 8, 2},
		{"start of line3", 12, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := line.DistanceFromLF(buf, tt.from)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFindLineEnd(t *testing.T) {
	buf := []rune("line1\nline2\nline3")

	tests := []struct {
		name  string
		start int
		want  int
	}{
		{"start of line1", 0, 5},
		{"start of line2", 6, 11},
		{"start of line3", 12, 17},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := line.FindLineEnd(buf, tt.start)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFindNextLineStart(t *testing.T) {
	buf := []rune("line1\nline2\nline3")

	tests := []struct {
		name string
		from int
		want int
	}{
		{"start of buffer", 0, 6},
		{"middle of line1", 3, 6},
		{"start of line2", 6, 12},
		{"end of buffer", 17, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := line.FindNextLineStart(buf, tt.from)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClampToLine(t *testing.T) {
	buf := []rune("line1\nline2\nline3")

	tests := []struct {
		name      string
		lineStart int
		col       int
		want      int
	}{
		{"column within line", 0, 3, 3},
		{"column at end of line", 0, 5, 5},
		{"column past line", 0, 10, 5},
		{"line2, column within", 6, 2, 8},
		{"line3, column past line", 12, 10, 17},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := line.ClampToLine(buf, tt.lineStart, tt.col)
			assert.Equal(t, tt.want, got)
		})
	}
}
