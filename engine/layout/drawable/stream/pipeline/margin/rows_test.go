package margin

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestRowsTransformer_KeepHasNext(t *testing.T) {
	transformer := ColsCenter(2)

	mock := drawable_test.MockUnit{}

	lines := []text.Line{
		*text.NewLine("golang"),
		*text.NewLine("ziglang"),
	}

	size := winsize.Winsize{
		Rows: 10,
		Cols: 10,
	}

	result, hasNext := transformer(
		size, mock.ToUnit(), lines, true,
	)

	assert.True(t, hasNext)
	assert.NotNil(t, result)
}

func TestRowsTopTransformer(t *testing.T) {
	mock := drawable_test.MockUnit{}

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Rows
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(10, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   3,
			wantLines: []string{"", "", "golang"},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(2, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   2,
			wantLines: []string{"", "golang"},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(1, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := RowsTop(tt.margin)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result, _ := transformer(
				tt.size, mock.ToUnit(), lines, true,
			)

			assert.Len(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, text.LineToString(&result[i]))
			}
		})
	}
}

func TestRowsBottomTransformer(t *testing.T) {
	mock := drawable_test.MockUnit{}

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Rows
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(10, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   3,
			wantLines: []string{"golang", "", ""},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(2, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   2,
			wantLines: []string{"golang", ""},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(1, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := RowsBottom(tt.margin)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result, _ := transformer(
				tt.size, mock.ToUnit(), lines, true,
			)

			assert.Len(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, text.LineToString(&result[i]))
			}
		})
	}
}

func TestRowsMiddleTransformer(t *testing.T) {
	mock := drawable_test.MockUnit{}

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Rows
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(10, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   5,
			wantLines: []string{"", "", "golang", "", ""},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(3, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   3,
			wantLines: []string{"", "golang", ""},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(1, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := RowsMiddle(tt.margin)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result, _ := transformer(
				tt.size, mock.ToUnit(), lines, true,
			)

			assert.Len(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, text.LineToString(&result[i]))
			}
		})
	}
}
