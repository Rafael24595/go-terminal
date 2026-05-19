package padding

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func renderFragments(styler *styler.Spec, size winsize.Winsize, frags []text.Fragment) string {
	var buffer strings.Builder

	lineSize := winsize.New(
		size.Rows,
		size.Cols,
	)

	for _, f := range frags {
		spec := styler.Apply(f.Spec, lineSize, f.Text, f.Size())

		fragSize := text.FragmentMeasure(size.Cols, f)
		lineSize.Cols = lineSize.Cols.Sub(fragSize)

		buffer.WriteString(spec)
	}

	return buffer.String()
}

func TestColPositioners(t *testing.T) {
	tests := []struct {
		name      string
		position  style.HorizontalPosition
		remaining winsize.Cols
		wantLeft  winsize.Cols
		wantRight winsize.Cols
	}{
		{
			name:      "Left alignment puts everything on the right",
			position:  style.Left,
			remaining: 10,
			wantLeft:  0,
			wantRight: 10,
		},
		{
			name:      "Right alignment puts everything on the left",
			position:  style.Right,
			remaining: 10,
			wantLeft:  10,
			wantRight: 0,
		},
		{
			name:      "Center alignment with even number splits equally",
			position:  style.Center,
			remaining: 6,
			wantLeft:  3,
			wantRight: 3,
		},
		{
			name:      "Center alignment with odd number handles the remainder",
			position:  style.Center,
			remaining: 5,
			wantLeft:  2,
			wantRight: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positioner := colPositionerMap[tt.position]
			gotLeft, gotRight := positioner(tt.remaining)

			assert.Equal(t, tt.wantLeft, gotLeft)
			assert.Equal(t, tt.wantRight, gotRight)
		})
	}
}

func TestColsTransformer(t *testing.T) {
	styler := styler.NewDefault()

	mockLines := []text.Line{
		*text.NewLine("Golang"),
	}

	tests := []struct {
		name       string
		hint       SizeHint[winsize.Cols]
		pos        style.HorizontalPosition
		size       winsize.Winsize
		wantLength winsize.Cols
		wantString string
	}{
		{
			name:       "Fixed size adds padding left",
			hint:       Fixed(winsize.Cols(10)),
			pos:        style.Left,
			size:       winsize.New(10, 20),
			wantLength: 10,
			wantString: "Golang    ",
		},
		{
			name:       "Fixed size adds padding right",
			hint:       Fixed(winsize.Cols(10)),
			pos:        style.Right,
			size:       winsize.New(10, 20),
			wantLength: 10,
			wantString: "    Golang",
		},
		{
			name:       "Fixed size adds padding center",
			hint:       Fixed(winsize.Cols(10)),
			pos:        style.Center,
			size:       winsize.New(10, 20),
			wantLength: 10,
			wantString: "  Golang  ",
		},
		{
			name:       "Maximize size adds padding center",
			hint:       Maximize[winsize.Cols](),
			pos:        style.Center,
			size:       winsize.New(10, 12),
			wantLength: 12,
			wantString: "   Golang   ",
		},
		{
			name:       "Clamp zero ensures no panic if line is wider than min",
			hint:       Fixed(winsize.Cols(4)),
			pos:        style.Left,
			size:       winsize.New(10, 20),
			wantLength: 6,
			wantString: "Golang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Cols(tt.hint, tt.pos)
			resLines := transformer(tt.size, mockLines)

			gotSize := text.FragmentMeasure(tt.size.Cols, resLines[0].Text...)

			assert.Equal(t, tt.wantLength, gotSize)
			assert.Equal(t, tt.wantString, renderFragments(styler, tt.size, resLines[0].Text))
		})
	}
}
