package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestHStack_DrawableBasicSuite(t *testing.T) {
	dw := HStackDrawableFromDrawables()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestHStack_Distribution(t *testing.T) {
	d1 := &drawable_test.MockDrawable{}
	d2 := &drawable_test.MockDrawable{}
	d3 := &drawable_test.MockDrawable{}

	stack := NewHStackDrawable(
		d1.ToDrawable(),
		d2.ToDrawable(),
		d3.ToDrawable(),
	)

	stack.init()

	layers := stack.fixed

	total := uint16(0)
	for _, l := range layers {
		total += l.chunk
	}

	assert.Equal(t, 100, total)
	assert.Equal(t, 34, layers[0].chunk)
	assert.Equal(t, 33, layers[1].chunk)
}

func TestHStack_MixedFixedAndDynamic(t *testing.T) {
	d1 := &drawable_test.MockDrawable{}
	d2 := &drawable_test.MockDrawable{}
	d3 := &drawable_test.MockDrawable{}

	stack := NewHStackDrawable()
	stack.PushChunk(d1.ToDrawable(), 20)
	stack.Push(d2.ToDrawable(), d3.ToDrawable())

	stack.init()

	assert.Equal(t, 20, stack.fixed[0].chunk)
	assert.Equal(t, 40, stack.fixed[1].chunk)
	assert.Equal(t, 40, stack.fixed[2].chunk)
}

func TestHStack_RenderOutput(t *testing.T) {
	size := terminal.Winsize{Rows: 1, Cols: 6}

	dA := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("go-lang"),
		},
	}
	dB := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("ziglang"),
		},
	}

	stack := NewHStackDrawable()
	stack.PushChunk(dA.ToDrawable(), 50)
	stack.PushChunk(dB.ToDrawable(), 50)
	stack.init()

	lines, _ := stack.draw(size)

	assert.Len(t, 3, lines)

	resultText := ""
	for _, frag := range lines[0].Text {
		resultText += frag.Text
	}

	assert.Equal(t, "go-zig", text.LineToString(&lines[0]))
	assert.Equal(t, "lanlan", text.LineToString(&lines[1]))
	assert.Equal(t, "gg", text.LineToString(&lines[2]))
}
