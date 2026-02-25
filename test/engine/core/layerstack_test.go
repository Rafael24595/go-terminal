package core_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestLayerStack_Init(t *testing.T) {
	stack := &core.LayerStack{}

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

func TestLayerStack_Shift_Order(t *testing.T) {
	stack := &core.LayerStack{}

	count := 0

	m1 := &drawable_test.MockDrawable{Status: false}
	m2 := &drawable_test.MockDrawable{Status: false}

	d1 := m1.ToDrawable()
	d2 := m2.ToDrawable()

	d1.Draw = func() ([]core.Line, bool) {
		m1.Order = count
		count++
		return m1.Draw()
	}

	d2.Draw = func() ([]core.Line, bool) {
		m2.Order = count
		count++
		return m2.Draw()
	}

	stack.Shift(d1)
	stack.Shift(d2)

	stack.Draw()

	assert.Equal(t, 0, m1.Order)
	assert.Equal(t, 1, m2.Order)
}

func TestLayerStack_Unshift_Order(t *testing.T) {
	stack := &core.LayerStack{}

	count := 0

	m1 := &drawable_test.MockDrawable{Status: false}
	m2 := &drawable_test.MockDrawable{Status: false}

	d1 := m1.ToDrawable()
	d2 := m2.ToDrawable()

	d1.Draw = func() ([]core.Line, bool) {
		m1.Order = count
		count++
		return m1.Draw()
	}

	d2.Draw = func() ([]core.Line, bool) {
		m2.Order = count
		count++
		return m2.Draw()
	}

	stack.Shift(d1)
	stack.Unshift(d2)

	stack.Draw()

	assert.Equal(t, 1, m1.Order)
	assert.Equal(t, 0, m2.Order)
}

func TestLayerStack_Draw_BreaksOnTrue(t *testing.T) {
	stack := &core.LayerStack{}

	d1 := &drawable_test.MockDrawable{Status: true}
	d2 := &drawable_test.MockDrawable{Status: false}

	stack.Shift(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	_, global := stack.Draw()

	assert.True(t, global)
	assert.Equal(t, 0, d2.DrawCalls)
}

func TestLayerStack_DisablesLayer(t *testing.T) {
	stack := &core.LayerStack{}

	d1 := &drawable_test.MockDrawable{Status: false}

	stack.Shift(d1.ToDrawable())

	stack.Draw()
	stack.Draw()

	assert.Equal(t, 1, d1.DrawCalls)
}

func TestLayerStack_BufferConcat(t *testing.T) {
	stack := &core.LayerStack{}

	line1 := core.LineFromString("go")
	line2 := core.LineFromString("lang")

	d1 := &drawable_test.MockDrawable{
		Lines:  []core.Line{line1},
		Status: false,
	}

	d2 := &drawable_test.MockDrawable{
		Lines:  []core.Line{line2},
		Status: false,
	}

	stack.Shift(
		d1.ToDrawable(),
		d2.ToDrawable(),
	)

	buffer, _ := stack.Draw()

	assert.Len(t, 2, buffer)
	assert.Equal(t, "golang", core.LineToString(buffer[0])+core.LineToString(buffer[1]))
}

func TestLayerStack_ShortCircuitStopsPropagation(t *testing.T) {
	stack := &core.LayerStack{}

	d1 := &drawable_test.MockDrawable{Status: false}
	d2 := &drawable_test.MockDrawable{Status: true}
	d3 := &drawable_test.MockDrawable{Status: false}

	stack.Shift(
		d1.ToDrawable(),
		d2.ToDrawable(),
		d3.ToDrawable(),
	)

	stack.Draw()

	assert.Equal(t, 1, d1.DrawCalls)
	assert.Equal(t, 1, d2.DrawCalls)
	assert.Equal(t, 0, d3.DrawCalls)
}
