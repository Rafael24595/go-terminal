package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestHStack_DrawableBasicSuite(t *testing.T) {
	dw := HStackDrawableFromDrawables()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestHStack_Distribution(t *testing.T) {
	d1 := &drawable_test.MockDrawable{}
	d2 := &drawable_test.MockDrawable{}
	d3 := &drawable_test.MockDrawable{}

	stack := NewHStack(
		d1.ToDrawable(),
		d2.ToDrawable(),
		d3.ToDrawable(),
	)

	stack.init()
	stack.lazyInit(winsize.Winsize{
		Cols: 100,
	})

	assert.Len(t, 3, stack.fixed)

	assert.Equal(t, 34, stack.fixed[0].value)
	assert.Equal(t, 33, stack.fixed[1].value)
	assert.Equal(t, 33, stack.fixed[2].value)
}

func TestHStack_MixedFixedAndDynamic(t *testing.T) {
	d1 := &drawable_test.MockDrawable{}
	d2 := &drawable_test.MockDrawable{}
	d3 := &drawable_test.MockDrawable{}

	stack := NewHStack()
	stack.PushChunk(d1.ToDrawable(), chunk.Fixed[winsize.Cols](20))
	stack.Push(d2.ToDrawable(), d3.ToDrawable())

	stack.init()
	stack.lazyInit(winsize.Winsize{
		Cols: 100,
	})

	assert.Equal(t, 20, stack.fixed[0].value)
	assert.Equal(t, 40, stack.fixed[1].value)
	assert.Equal(t, 40, stack.fixed[2].value)
}

func TestHStack_RenderOutput(t *testing.T) {
	size := winsize.Winsize{Rows: 1, Cols: 6}

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

	stack := NewHStack()
	stack.PushChunk(dA.ToDrawable(), chunk.Percent[winsize.Cols](50))
	stack.PushChunk(dB.ToDrawable(), chunk.Percent[winsize.Cols](50))

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

func TestHStack_ToDrawable_AnemicStack(t *testing.T) {
    dw1 := &drawable_test.MockDrawable{
		Name: "mock_drawable",
	}

	stack := NewHStack().
		PushChunk(dw1.ToDrawable(), chunk.Dynamic[winsize.Cols]()).
		ToDrawable()

    assert.True(t, stack.Tags.Has(AnemicStack))
    assert.Equal(t, dw1.Name, stack.Name) 
}

func TestHStack_ToDrawable_SingleElement_NotAnemic(t *testing.T) {
    dw1 := &drawable_test.MockDrawable{
		Name: "mock_drawable",
	}

	stack := NewHStack().
		PushChunk(dw1.ToDrawable(), chunk.Fixed[winsize.Cols](10)).
		ToDrawable()

    assert.False(t, stack.Tags.Has(AnemicStack))
    assert.Equal(t, NameHStack, stack.Name) 
}

func TestHStack_ToDrawable_MultipleElements(t *testing.T) {
    dw1 := &drawable_test.MockDrawable{
		Name: "mock_drawable_001",
	}

	dw2 := &drawable_test.MockDrawable{
		Name: "mock_drawable_002",
	}

	stack := NewHStack().
		PushChunk(dw1.ToDrawable(), chunk.Dynamic[winsize.Cols]()).
		PushChunk(dw2.ToDrawable(), chunk.Dynamic[winsize.Cols]()).
		ToDrawable()

    assert.False(t, stack.Tags.Has(AnemicStack))
    assert.Equal(t, NameHStack, stack.Name) 
}
