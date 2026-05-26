package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestVStack_UnitBasicSuite(t *testing.T) {
	unit := VStackFromUnits()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestVStack_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	mock1 := &drawable_test.MockUnit{}

	unit := NewVStack().Push(
		mock1.ToUnit(),
	)

	unit.init()

	assert.Panic(t, func() {
		m2 := &drawable_test.MockUnit{}
		unit.Push(m2.ToUnit())
	})
}

func TestVStack_Init(t *testing.T) {
	stack := &VStackUnit{}

	mock1 := &drawable_test.MockUnit{}
	mock2 := &drawable_test.MockUnit{}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
	)

	stack.init()
	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Greater(t, 0, mock1.InitCalled)
	assert.Greater(t, 0, mock2.InitCalled)
}

func TestVStack_Shift_Order(t *testing.T) {
	stack := &VStackUnit{}

	count := uint(0)

	mock1 := &drawable_test.MockUnit{Status: false}
	mock2 := &drawable_test.MockUnit{Status: false}

	unit1 := mock1.ToUnit()
	unit2 := mock2.ToUnit()

	unit1.Drawable.Draw = func(_ winsize.Winsize) ([]text.Line, bool) {
		mock1.DrawCalls = count
		count++
		return make([]text.Line, 0), false
	}

	unit2.Drawable.Draw = func(_ winsize.Winsize) ([]text.Line, bool) {
		mock2.DrawCalls = count
		count++
		return make([]text.Line, 0), false
	}

	stack.Push(unit1)
	stack.Push(unit2)

	stack.init()

	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Equal(t, 0, mock1.DrawCalls)
	assert.Equal(t, 1, mock2.DrawCalls)
}

func TestVStack_Unshift_Order(t *testing.T) {
	stack := &VStackUnit{}

	count := uint(0)

	mock1 := &drawable_test.MockUnit{Status: false}
	mock2 := &drawable_test.MockUnit{Status: false}

	unit1 := mock1.ToUnit()
	unit2 := mock2.ToUnit()

	unit1.Drawable.Draw = func(_ winsize.Winsize) ([]text.Line, bool) {
		mock1.DrawCalls = count
		count++
		return make([]text.Line, 0), false
	}

	unit2.Drawable.Draw = func(_ winsize.Winsize) ([]text.Line, bool) {
		mock2.DrawCalls = count
		count++
		return make([]text.Line, 0), false
	}

	stack.Push(unit1)
	stack.Unshift(unit2)

	stack.init()

	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Equal(t, 1, mock1.DrawCalls)
	assert.Equal(t, 0, mock2.DrawCalls)
}

func TestVStack_Draw_BreaksOnTrue(t *testing.T) {
	stack := &VStackUnit{}

	mock1 := &drawable_test.MockUnit{Status: true}
	mock2 := &drawable_test.MockUnit{Status: false}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
	)

	stack.init()

	_, global := stack.draw(winsize.Winsize{})

	assert.True(t, global)
	assert.Equal(t, 0, mock2.DrawCalls)
}

func TestVStack_DisablesLayer(t *testing.T) {
	stack := &VStackUnit{}

	mock := &drawable_test.MockUnit{Status: false}

	stack.Push(mock.ToUnit())

	stack.init()

	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})
	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Equal(t, 1, mock.DrawCalls)
}

func TestVStack_BufferConcat(t *testing.T) {
	stack := &VStackUnit{}

	line1 := text.NewLine("go")
	line2 := text.NewLine("lang")

	mock1 := &drawable_test.MockUnit{
		Lines:  []text.Line{*line1},
		Status: false,
	}
	mock2 := &drawable_test.MockUnit{
		Lines:  []text.Line{*line2},
		Status: false,
	}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
	)

	stack.init()

	buffer, _ := stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Len(t, 2, buffer)
	assert.Equal(t, "golang", text.LineToString(&buffer[0])+text.LineToString(&buffer[1]))
}

func TestVStack_ShortCircuitStopsPropagation(t *testing.T) {
	stack := &VStackUnit{}

	mock1 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 1),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 2),
	}
	mock3 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 1),
	}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
		mock3.ToUnit(),
	)

	stack.init()

	stack.draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.Equal(t, 1, mock1.DrawCalls)
	assert.Equal(t, 1, mock2.DrawCalls)
	assert.Equal(t, 0, mock3.DrawCalls)
}

func TestVStack_FixedChunk_PadsWhenChildIsSmaller(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: make([]text.Line, 10),
	}

	stack := NewVStack().
		PushLayer(
			mock.ToUnit(),
			layer.Fixed[winsize.Rows](15),
		).
		ToUnit()

	stack.Drawable.Init()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 20, Cols: 10})

	assert.Len(t, 15, lines)
}

func TestVStack_FixedChunk_TruncatesWhenChildIsBigger(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: make([]text.Line, 15),
	}

	stack := NewVStack().
		PushLayer(
			mock.ToUnit(),
			layer.Fixed[winsize.Rows](20),
		).
		ToUnit()

	stack.Drawable.Init()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 10, Cols: 10})

	assert.Len(t, 10, lines)
}

func TestVStack_DynamicChunk_FillsRemainingSpace(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 10),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 10),
	}
	mock3 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 5),
	}

	stack := NewVStack().
		PushLayer(
			mock1.ToUnit(),
			layer.Fixed[winsize.Rows](10),
		).
		PushLayer(mock2.ToUnit()).
		PushLayer(mock3.ToUnit()).
		ToUnit()

	stack.Drawable.Init()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 30, Cols: 10})

	assert.Len(t, 25, lines)
}

func TestVStack_FixedOverflow_ShouldNotExceedContainer(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 10),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 10),
	}

	stack := NewVStack().
		PushLayer(
			mock1.ToUnit(),
			layer.Fixed[winsize.Rows](10),
		).
		PushLayer(
			mock2.ToUnit(),
			layer.Fixed[winsize.Rows](10),
		).
		ToUnit()

	stack.Drawable.Init()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 15, Cols: 10})

	assert.Len(t, 15, lines)
}

func TestVStack_ExactFit_NoExtraNoMissing(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 5),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 5),
	}
	mock3 := &drawable_test.MockUnit{
		Lines: make([]text.Line, 5),
	}

	stack := NewVStack().
		PushLayer(mock1.ToUnit()).
		PushLayer(mock2.ToUnit()).
		PushLayer(mock3.ToUnit()).
		ToUnit()

	stack.Drawable.Init()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 15, Cols: 10})

	assert.Len(t, 15, lines)
}

func TestVStack_ToUnit_AnemicStack(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewVStack().
		PushLayer(mock.ToUnit()).
		ToUnit()

	assert.True(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, mock.Name, stack.Name)
}

func TestVStack_ToUnit_SingleElement_NotAnemic(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewVStack().
		PushLayer(mock.ToUnit(), layer.Fixed[winsize.Rows](10)).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameVStack, stack.Name)
}

func TestVStack_ToUnit_MultipleElements(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Name: "mock_unit_001",
	}
	mock2 := &drawable_test.MockUnit{
		Name: "mock_unit_002",
	}

	stack := NewVStack().
		PushLayer(mock1.ToUnit()).
		PushLayer(mock2.ToUnit()).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameVStack, stack.Name)
}
