package padding

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestRowPositioners(t *testing.T) {
	inputLines := []text.Line{
		*text.NewLine("Golang"),
	}

	tests := []struct {
		name         string
		position     style.VerticalPosition
		paddingTotal winsize.Rows
		wantLength   int
		wantCursor   int
	}{
		{
			name:         "Top alignment stretches to padding total",
			position:     style.Top,
			paddingTotal: 5,
			wantLength:   5,
			wantCursor:   0,
		},
		{
			name:         "Bottom alignment stretches to padding total",
			position:     style.Bottom,
			paddingTotal: 5,
			wantLength:   5,
			wantCursor:   4,
		},
		{
			name:         "Center alignment stretches to padding total",
			position:     style.Middle,
			paddingTotal: 5,
			wantLength:   5,
			wantCursor:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positioner := rowPositionerMap[tt.position]
			gotLines := positioner(inputLines, tt.paddingTotal)

			assert.Len(t, tt.wantLength, gotLines)

			for i := range gotLines {
				wantString := ""
				if i == tt.wantCursor {
					wantString = "Golang"
				}

				assert.Equal(t, wantString, text.LineToString(&gotLines[i]))
			}
		})
	}
}

func TestRowsTransformer(t *testing.T) {
	mockLines := []text.Line{
		*text.NewLine("Golang"),
	}

	tests := []struct {
		name       string
		hint       hint.Size[winsize.Rows]
		pos        style.VerticalPosition
		canvasSize winsize.Winsize
		wantRows   int
	}{
		{
			name:       "Fixed size expands rows if input is shorter",
			hint:       hint.Fixed(winsize.Rows(4)),
			pos:        style.Top,
			canvasSize: winsize.New(10, 20),
			wantRows:   4,
		},
		{
			name:       "Guard clause keeps lines intact if already larger than min",
			hint:       hint.Fixed(winsize.Rows(1)),
			pos:        style.Top,
			canvasSize: winsize.New(10, 20),
			wantRows:   1,
		},
		{
			name:       "Maximize hint takes all available height from canvas",
			hint:       hint.Maximize[winsize.Rows](),
			pos:        style.Bottom,
			canvasSize: winsize.New(7, 20),
			wantRows:   7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Rows(tt.hint, tt.pos)

			gotLines := transformer(tt.canvasSize, mockLines)

			assert.Len(t, tt.wantRows, gotLines)
		})
	}
}
