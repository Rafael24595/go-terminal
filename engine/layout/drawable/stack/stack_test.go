package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestStackDrawable_ToDrawable(t *testing.T) {
	dw := StackDrawableFromDrawables()
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestStackDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	bd := NewStackDrawable()

	assert.Panic(t, func() {
		bd.Draw()
	})
}

func TestStackDrawable_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	bd := NewStackDrawable()

	m1 := &drawable_test.MockDrawable{}
	bd.Shift(m1.ToDrawable())

	bd.Init(terminal.Winsize{})

	assert.Panic(t, func() {
		m2 := &drawable_test.MockDrawable{}
		bd.Shift(m2.ToDrawable())
	})
}

func TestStackDrawable_Init(t *testing.T) {
	stack := &StackDrawable{}

	d1 := &drawable_test.MockDrawable{}
	d2 := &drawable_test.MockDrawable{}

	stack.Shift(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	stack.Init(terminal.Winsize{})

	assert.True(t, d1.InitCalled)
	assert.True(t, d2.InitCalled)
}

func TestStackDrawable_Shift_Order(t *testing.T) {
	stack := &StackDrawable{}

	count := 0

	m1 := &drawable_test.MockDrawable{Status: false}
	m2 := &drawable_test.MockDrawable{Status: false}

	d1 := m1.ToDrawable()
	d2 := m2.ToDrawable()

	d1.Draw = func() ([]text.Line, bool) {
		m1.Order = count
		count++
		return m1.Draw()
	}

	d2.Draw = func() ([]text.Line, bool) {
		m2.Order = count
		count++
		return m2.Draw()
	}

	stack.Shift(d1)
	stack.Shift(d2)

	stack.Init(terminal.Winsize{})

	stack.Draw()

	assert.Equal(t, 0, m1.Order)
	assert.Equal(t, 1, m2.Order)
}

func TestStackDrawable_Unshift_Order(t *testing.T) {
	stack := &StackDrawable{}

	count := 0

	m1 := &drawable_test.MockDrawable{Status: false}
	m2 := &drawable_test.MockDrawable{Status: false}

	d1 := m1.ToDrawable()
	d2 := m2.ToDrawable()

	d1.Draw = func() ([]text.Line, bool) {
		m1.Order = count
		count++
		return m1.Draw()
	}

	d2.Draw = func() ([]text.Line, bool) {
		m2.Order = count
		count++
		return m2.Draw()
	}

	stack.Shift(d1)
	stack.Unshift(d2)

	stack.Init(terminal.Winsize{})

	stack.Draw()

	assert.Equal(t, 1, m1.Order)
	assert.Equal(t, 0, m2.Order)
}

func TestStackDrawable_Draw_BreaksOnTrue(t *testing.T) {
	stack := &StackDrawable{}

	d1 := &drawable_test.MockDrawable{Status: true}
	d2 := &drawable_test.MockDrawable{Status: false}

	stack.Shift(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	stack.Init(terminal.Winsize{})

	_, global := stack.Draw()

	assert.True(t, global)
	assert.Equal(t, 0, d2.DrawCalls)
}

func TestStackDrawable_DisablesLayer(t *testing.T) {
	stack := &StackDrawable{}

	d1 := &drawable_test.MockDrawable{Status: false}

	stack.Shift(d1.ToDrawable())

	stack.Init(terminal.Winsize{})

	stack.Draw()
	stack.Draw()

	assert.Equal(t, 1, d1.DrawCalls)
}

func TestStackDrawable_BufferConcat(t *testing.T) {
	stack := &StackDrawable{}

	line1 := text.LineFromString("go")
	line2 := text.LineFromString("lang")

	d1 := &drawable_test.MockDrawable{
		Lines:  []text.Line{line1},
		Status: false,
	}

	d2 := &drawable_test.MockDrawable{
		Lines:  []text.Line{line2},
		Status: false,
	}

	stack.Shift(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	stack.Init(terminal.Winsize{})

	buffer, _ := stack.Draw()

	assert.Len(t, 2, buffer)
	assert.Equal(t, "golang", text.LineToString(buffer[0])+text.LineToString(buffer[1]))
}

func TestStackDrawable_ShortCircuitStopsPropagation(t *testing.T) {
	stack := &StackDrawable{}

	d1 := &drawable_test.MockDrawable{Status: false}
	d2 := &drawable_test.MockDrawable{Status: true}
	d3 := &drawable_test.MockDrawable{Status: false}

	stack.Shift(
		d1.ToDrawable(),
		d2.ToDrawable(),
		d3.ToDrawable(),
	)

	stack.Init(terminal.Winsize{})

	stack.Draw()

	assert.Equal(t, 1, d1.DrawCalls)
	assert.Equal(t, 1, d2.DrawCalls)
	assert.Equal(t, 0, d3.DrawCalls)
}
