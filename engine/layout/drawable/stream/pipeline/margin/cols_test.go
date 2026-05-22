package margin

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	render_test "github.com/Rafael24595/go-reacterm-core/test/engine/render"
)

func TestColsTransformer_KeepHasNext(t *testing.T) {
	transformer := ColsCenter(2)

	mock := drawable_test.MockUnit{}

	lines := []text.Line{
		*text.NewLine("golang"),
	}

	size := winsize.Winsize{
		Rows: 5,
		Cols: 10,
	}

	result, hasNext := transformer(
		size, mock.ToUnit(), lines, true,
	)

	assert.True(t, hasNext)
	assert.NotNil(t, result)
}

func TestColsLeftTransformer(t *testing.T) {
	mock := drawable_test.MockUnit{}
	styler := styler.NewDefault()

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Cols
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(5, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"  golang"},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(5, 7),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{" golang"},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(5, 6),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := ColsLeft(tt.margin)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result, _ := transformer(
				tt.size, mock.ToUnit(), lines, true,
			)

			assert.Len(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, render_test.Fragments(styler, tt.size, result[i].Text))
			}
		})
	}
}

func TestColsRightTransformer(t *testing.T) {
	mock := drawable_test.MockUnit{}
	styler := styler.NewDefault()

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Cols
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(5, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang  "},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(5, 7),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang "},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(5, 6),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := ColsRight(tt.margin)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result, _ := transformer(
				tt.size, mock.ToUnit(), lines, true,
			)

			assert.Len(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, render_test.Fragments(styler, tt.size, result[i].Text))
			}
		})
	}
}

func TestColsCenterTransformer(t *testing.T) {
	mock := drawable_test.MockUnit{}
	styler := styler.NewDefault()

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Cols
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(5, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"  golang  "},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(5, 8),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{" golang "},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(5, 6),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := ColsCenter(tt.margin)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result, _ := transformer(
				tt.size, mock.ToUnit(), lines, true,
			)

			assert.Len(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, render_test.Fragments(styler, tt.size, result[i].Text))
			}
		})
	}
}
