package drawable_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

type MockDrawable struct {
	Order      int
	InitCalled bool
	DrawCalls  int
	Lines      []core.Line
	Status     bool
}

func (m *MockDrawable) Init(size terminal.Winsize) {
	m.InitCalled = true
}

func (m *MockDrawable) Draw() ([]core.Line, bool) {
	m.DrawCalls++
	return m.Lines, m.Status
}

func (m *MockDrawable) ToDrawable() core.Drawable {
	return core.Drawable{
		Init: m.Init,
		Draw: m.Draw,
	}
}

func Helper_ToDrawable(t *testing.T, drawable core.Drawable) {
	t.Helper()

	assert.NotNil(t, drawable.Init, "Drawable.Init should be set")
	assert.NotNil(t, drawable.Draw, "Drawable.Draw should be set")
}
