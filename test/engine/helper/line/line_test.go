package line_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/line"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
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

	index, _ := line.FindNextLineStart(buf, 0)
	assert.Equal(t, 4, index)

	index, _ = line.FindNextLineStart(buf, 2)
	assert.Equal(t, 4, index)

	index, ok := line.FindNextLineStart(buf, 4)
	assert.False(t, ok)
	assert.Equal(t, 0, index)
}

func TestFindPrevLineStart_EmptyLine(t *testing.T) {
	buf := []rune("abc\n\ndef")
	got, _ := line.FindPrevLineStart(buf, 5)
	assert.Equal(t, 4, got)
}

func TestFindPrevLineStart_MultipleEmptyLines(t *testing.T) {
	buf := []rune("a\n\n\nb")
	got, _ := line.FindPrevLineStart(buf, 4)
	assert.Equal(t, 3, got)
}

func TestFindPrevLineStart_Normal(t *testing.T) {
	buf := []rune("a\nb\nc")
	got, _ := line.FindPrevLineStart(buf, 4)
	assert.Equal(t, 2, got)
}

func TestClampToLine_EmptyLine(t *testing.T) {
	buf := []rune("abc\n\ndef")
	got := line.ClampToLine(buf, 4, 10)
	assert.Equal(t, 4, got)
}

func TestFindLineStart(t *testing.T) {
	buf := []rune("line1\nline2\nline3")

	tests := []struct {
		name string
		from offset.Offset
		want offset.Offset
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
		from offset.Offset
		want offset.Offset
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
		start offset.Offset
		want  offset.Offset
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
		from offset.Offset
		want offset.Offset
		stat bool
	}{
		{"start of buffer", 0, 6, true},
		{"middle of line1", 3, 6, true},
		{"start of line2", 6, 12, true},
		{"end of buffer", 17, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := line.FindNextLineStart(buf, tt.from)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.stat, ok)
		})
	}
}

func TestClampToLine(t *testing.T) {
	buf := []rune("line1\nline2\nline3")

	tests := []struct {
		name      string
		lineStart offset.Offset
		col       offset.Offset
		want      offset.Offset
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
