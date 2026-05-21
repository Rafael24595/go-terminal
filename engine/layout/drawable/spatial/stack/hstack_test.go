package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestHStack_UnitBasicSuite(t *testing.T) {
	unit := HStackFromUnits()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestHStack_Distribution(t *testing.T) {
	mock1 := &drawable_test.MockUnit{}
	mock2 := &drawable_test.MockUnit{}
	mock3 := &drawable_test.MockUnit{}

	stack := NewHStack(
		mock1.ToUnit(),
		mock2.ToUnit(),
		mock3.ToUnit(),
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
	mock1 := &drawable_test.MockUnit{}
	mock2 := &drawable_test.MockUnit{}
	mock3 := &drawable_test.MockUnit{}

	stack := NewHStack()
	stack.PushChunk(mock1.ToUnit(), chunk.Fixed[winsize.Cols](20))
	stack.Push(mock2.ToUnit(), mock3.ToUnit())

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

	mock1 := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("go-lang"),
		},
	}
	mock2 := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("ziglang"),
		},
	}

	stack := NewHStack()
	stack.PushChunk(mock1.ToUnit(), chunk.Percent[winsize.Cols](50))
	stack.PushChunk(mock2.ToUnit(), chunk.Percent[winsize.Cols](50))

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

func TestHStack_ToUnit_AnemicStack(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewHStack().
		PushChunk(mock.ToUnit(), chunk.Dynamic[winsize.Cols]()).
		ToUnit()

	assert.True(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, mock.Name, stack.Name)
}

func TestHStack_ToUnit_SingleElement_NotAnemic(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewHStack().
		PushChunk(mock.ToUnit(), chunk.Fixed[winsize.Cols](10)).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameHStack, stack.Name)
}

func TestHStack_ToUnit_MultipleElements(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Name: "mock_unit_001",
	}
	mock2 := &drawable_test.MockUnit{
		Name: "mock_unit_002",
	}

	stack := NewHStack().
		PushChunk(mock1.ToUnit(), chunk.Dynamic[winsize.Cols]()).
		PushChunk(mock2.ToUnit(), chunk.Dynamic[winsize.Cols]()).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameHStack, stack.Name)
}
