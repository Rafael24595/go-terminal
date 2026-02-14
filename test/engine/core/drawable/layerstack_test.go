package drawable_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

type mockDrawable struct {
	order      int
	initCalled bool
	drawCalls  int
	lines      []core.Line
	status     bool
}

func (m *mockDrawable) Init(size terminal.Winsize) {
	m.initCalled = true
}

func (m *mockDrawable) Draw() ([]core.Line, bool) {
	m.drawCalls++
	return m.lines, m.status
}

func (m *mockDrawable) toDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: m.Init,
		Draw: m.Draw,
	}
}

func TestLayerStack_Init(t *testing.T) {
	stack := &drawable.LayerStack{}

	d1 := &mockDrawable{}
	d2 := &mockDrawable{}

	stack.Shift(
		d1.toDrawable(),
		d2.toDrawable(),
	)

	stack.Init(terminal.Winsize{})

	assert.True(t, d1.initCalled)
	assert.True(t, d2.initCalled)
}

func TestLayerStack_Shift_Order(t *testing.T) {
	stack := &drawable.LayerStack{}

	count := 0

	m1 := &mockDrawable{status: false}
	m2 := &mockDrawable{status: false}

	d1 := m1.toDrawable()
	d2 := m2.toDrawable()

	d1.Draw = func() ([]core.Line, bool) {
		m1.order = count
		count++
		return m1.Draw()
	}

	d2.Draw = func() ([]core.Line, bool) {
		m2.order = count
		count++
		return m2.Draw()
	}

	stack.Shift(d1)
	stack.Shift(d2)

	stack.Draw()

	assert.Equal(t, 0, m1.order)
	assert.Equal(t, 1, m2.order)
}

func TestLayerStack_Unshift_Order(t *testing.T) {
	stack := &drawable.LayerStack{}

	count := 0

	m1 := &mockDrawable{status: false}
	m2 := &mockDrawable{status: false}

	d1 := m1.toDrawable()
	d2 := m2.toDrawable()

	d1.Draw = func() ([]core.Line, bool) {
		m1.order = count
		count++
		return m1.Draw()
	}

	d2.Draw = func() ([]core.Line, bool) {
		m2.order = count
		count++
		return m2.Draw()
	}

	stack.Shift(d1)
	stack.Unshift(d2)

	stack.Draw()

	assert.Equal(t, 1, m1.order)
	assert.Equal(t, 0, m2.order)
}

func TestLayerStack_Draw_BreaksOnTrue(t *testing.T) {
	stack := &drawable.LayerStack{}

	d1 := &mockDrawable{status: true}
	d2 := &mockDrawable{status: false}

	stack.Shift(
		d1.toDrawable(),
		d2.toDrawable(),
	)

	_, global := stack.Draw()

	assert.True(t, global)
	assert.Equal(t, 0, d2.drawCalls)
}

func TestLayerStack_DisablesLayer(t *testing.T) {
	stack := &drawable.LayerStack{}

	d1 := &mockDrawable{status: false}

	stack.Shift(d1.toDrawable())

	stack.Draw()
	stack.Draw()

	assert.Equal(t, 1, d1.drawCalls)
}

func TestLayerStack_BufferConcat(t *testing.T) {
	stack := &drawable.LayerStack{}

	line1 := core.LineFromString("go")
	line2 := core.LineFromString("lang")

	d1 := &mockDrawable{
		lines:  []core.Line{line1},
		status: false,
	}

	d2 := &mockDrawable{
		lines:  []core.Line{line2},
		status: false,
	}

	stack.Shift(
		d1.toDrawable(),
		d2.toDrawable(),
	)

	buffer, _ := stack.Draw()

	assert.Len(t, 2, buffer)
	assert.Equal(t, "golang", buffer[0].String() + buffer[1].String())
}

func TestLayerStack_ShortCircuitStopsPropagation(t *testing.T) {
	stack := &drawable.LayerStack{}

	d1 := &mockDrawable{status: false}
	d2 := &mockDrawable{status: true}
	d3 := &mockDrawable{status: false}

	stack.Shift(
		d1.toDrawable(),
		d2.toDrawable(),
		d3.toDrawable(),
	)

	stack.Draw()

	assert.Equal(t, 1, d1.drawCalls)
	assert.Equal(t, 1, d2.drawCalls)
	assert.Equal(t, 0, d3.drawCalls)
}
