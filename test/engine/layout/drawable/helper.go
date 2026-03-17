package drawable_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

type MockDrawable struct {
	Order      int
	InitCalled bool
	DrawCalls  int
	Lines      []text.Line
	Status     bool
}

func (m *MockDrawable) Init(size terminal.Winsize) {
	m.InitCalled = true
}

func (m *MockDrawable) Draw() ([]text.Line, bool) {
	m.DrawCalls++
	return m.Lines, m.Status
}

func (m *MockDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Init: m.Init,
		Draw: m.Draw,
	}
}

func Helper_ToDrawable(t *testing.T, drawable drawable.Drawable) {
	t.Helper()

	assert.NotEqual(t, "", drawable.Name, "Drawable.Name should be set")
	assert.NotNil(t, drawable.Code, "Drawable.Code should be set")
	assert.True(t, len(drawable.Tags) >= 0, "Drawable.Tags should be set")
	assert.NotNil(t, drawable.Init, "Drawable.Init should be set")
	assert.NotNil(t, drawable.Draw, "Drawable.Draw should be set")
}
