package drawable_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type MockDrawable struct {
	Order      int
	InitCalled bool
	WipeCalled bool
	DrawCalls  int
	Lines      []text.Line
	Status     bool
}

func (m *MockDrawable) Init() {
	m.InitCalled = true
}

func (m *MockDrawable) Wipe() {
	m.WipeCalled = true
}

func (m *MockDrawable) Draw(size terminal.Winsize) ([]text.Line, bool) {
	m.DrawCalls++
	return m.Lines, m.Status
}

func (m *MockDrawable) ToDrawable() drawable.Drawable {
	return drawable.Drawable{
		Name: "",
		Code: "",
		Tags: make(set.Set[string]),
		Init: m.Init,
		Wipe: m.Wipe,
		Draw: m.Draw,
	}
}

func Test_DrawableBasicSuite(t *testing.T, dw drawable.Drawable) {
	t.Helper()

	Helper_ToDrawable(t, dw)
	assert.Panic(t, func() {
		dw.Draw(terminal.Winsize{})
	})
}

func Helper_ToDrawable(t *testing.T, drawable drawable.Drawable) {
	t.Helper()

	assert.NotEqual(t, "", drawable.Name, "Drawable.Name should be set")
	assert.NotNil(t, drawable.Code, "Drawable.Code should be set")
	assert.True(t, len(drawable.Tags) >= 0, "Drawable.Tags should be set")
	assert.NotNil(t, drawable.Init, "Drawable.Init should be set")
	assert.NotNil(t, drawable.Wipe, "Drawable.Wipe should be set")
	assert.NotNil(t, drawable.Draw, "Drawable.Draw should be set")
}

func Helper_IsDrawableNil(t *testing.T, drawable drawable.Drawable) bool {
	t.Helper()

	return drawable.Init == nil &&
		drawable.Wipe == nil &&
		drawable.Draw == nil
}
