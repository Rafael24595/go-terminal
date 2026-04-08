package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestVStack_DrawableBasicSuite(t *testing.T) {
	dw := VStackDrawableFromDrawables()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestVStack_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	bd := NewVStackDrawable()

	m1 := &drawable_test.MockDrawable{}
	bd.Push(m1.ToDrawable())

	bd.init()

	assert.Panic(t, func() {
		m2 := &drawable_test.MockDrawable{}
		bd.Push(m2.ToDrawable())
	})
}

func TestVStack_Init(t *testing.T) {
	stack := &VStackDrawable{}

	d1 := &drawable_test.MockDrawable{}
	d2 := &drawable_test.MockDrawable{}

	stack.Push(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	stack.init()

	assert.True(t, d1.InitCalled)
	assert.True(t, d2.InitCalled)
}

func TestVStack_Shift_Order(t *testing.T) {
	stack := &VStackDrawable{}

	count := 0

	m1 := &drawable_test.MockDrawable{Status: false}
	m2 := &drawable_test.MockDrawable{Status: false}

	d1 := m1.ToDrawable()
	d2 := m2.ToDrawable()

	d1.Draw = func(_ terminal.Winsize) ([]text.Line, bool) {
		m1.Order = count
		count++
		return m1.Draw(terminal.Winsize{})
	}

	d2.Draw = func(_ terminal.Winsize) ([]text.Line, bool) {
		m2.Order = count
		count++
		return m2.Draw(terminal.Winsize{})
	}

	stack.Push(d1)
	stack.Push(d2)

	stack.init()

	stack.draw(terminal.Winsize{})

	assert.Equal(t, 0, m1.Order)
	assert.Equal(t, 1, m2.Order)
}

func TestVStack_Unshift_Order(t *testing.T) {
	stack := &VStackDrawable{}

	count := 0

	m1 := &drawable_test.MockDrawable{Status: false}
	m2 := &drawable_test.MockDrawable{Status: false}

	d1 := m1.ToDrawable()
	d2 := m2.ToDrawable()

	d1.Draw = func(_ terminal.Winsize) ([]text.Line, bool) {
		m1.Order = count
		count++
		return m1.Draw(terminal.Winsize{})
	}

	d2.Draw = func(_ terminal.Winsize) ([]text.Line, bool) {
		m2.Order = count
		count++
		return m2.Draw(terminal.Winsize{})
	}

	stack.Push(d1)
	stack.Unshift(d2)

	stack.init()

	stack.draw(terminal.Winsize{})

	assert.Equal(t, 1, m1.Order)
	assert.Equal(t, 0, m2.Order)
}

func TestVStack_Draw_BreaksOnTrue(t *testing.T) {
	stack := &VStackDrawable{}

	d1 := &drawable_test.MockDrawable{Status: true}
	d2 := &drawable_test.MockDrawable{Status: false}

	stack.Push(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	stack.init()

	_, global := stack.draw(terminal.Winsize{})

	assert.True(t, global)
	assert.Equal(t, 0, d2.DrawCalls)
}

func TestVStack_DisablesLayer(t *testing.T) {
	stack := &VStackDrawable{}

	d1 := &drawable_test.MockDrawable{Status: false}

	stack.Push(d1.ToDrawable())

	stack.init()

	stack.draw(terminal.Winsize{})
	stack.draw(terminal.Winsize{})

	assert.Equal(t, 1, d1.DrawCalls)
}

func TestVStack_BufferConcat(t *testing.T) {
	stack := &VStackDrawable{}

	line1 := text.NewLine("go")
	line2 := text.NewLine("lang")

	d1 := &drawable_test.MockDrawable{
		Lines:  []text.Line{line1},
		Status: false,
	}

	d2 := &drawable_test.MockDrawable{
		Lines:  []text.Line{line2},
		Status: false,
	}

	stack.Push(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	stack.init()

	buffer, _ := stack.draw(terminal.Winsize{})

	assert.Len(t, 2, buffer)
	assert.Equal(t, "golang", text.LineToString(buffer[0])+text.LineToString(buffer[1]))
}

func TestVStack_ShortCircuitStopsPropagation(t *testing.T) {
	stack := &VStackDrawable{}

	d1 := &drawable_test.MockDrawable{Status: false}
	d2 := &drawable_test.MockDrawable{Status: true}
	d3 := &drawable_test.MockDrawable{Status: false}

	stack.Push(
		d1.ToDrawable(),
		d2.ToDrawable(),
		d3.ToDrawable(),
	)

	stack.init()

	stack.draw(terminal.Winsize{})

	assert.Equal(t, 1, d1.DrawCalls)
	assert.Equal(t, 1, d2.DrawCalls)
	assert.Equal(t, 0, d3.DrawCalls)
}
