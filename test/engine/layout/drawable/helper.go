package drawable_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameMockUnit = "mock_unit"

//TODO: Refactor
type MockUnit struct {
	Order      int
	Tags       set.Set[string]
	Name       string
	InitCalled bool
	WipeCalled bool
	DrawCalls  int
	Lines      []text.Line
	queue      []text.Line
	Batch      uint
	Status     bool
	Size       winsize.Winsize
}

func (m *MockUnit) Init() {
	m.InitCalled = true
}

func (m *MockUnit) Wipe() {
	m.WipeCalled = true
}

func (m *MockUnit) Draw(size winsize.Winsize) ([]text.Line, bool) {
	m.DrawCalls++
	m.Size = size

	if m.Batch == 0 {
		return m.Lines, m.Status
	}

	if len(m.queue) == 0 {
		m.queue = m.Lines
	}

	limit := min(int(m.Batch), len(m.queue))

	data := m.queue[:limit]
	m.queue = m.queue[limit:]

	return data, len(m.queue) > 0
}

func (m *MockUnit) ToUnit() drawable.Unit {
	name := NameMockUnit
	if m.Name != "" {
		name = m.Name
	}

	return drawable.NewBuilder().
		Name(name).
		MergeTags(m.Tags).
		Init(m.Init).
		Wipe(m.Wipe).
		Draw(m.Draw).
		ToUnit()
}

func Test_UnitBasicSuite(t *testing.T, unit drawable.Unit) {
	t.Helper()

	Helper_ToUnit(t, unit)
	assert.Panic(t, func() {
		unit.Drawable.Draw(winsize.Winsize{})
	})
}

func Helper_ToUnit(t *testing.T, unit drawable.Unit) {
	t.Helper()

	assert.NotEqual(t, "", unit.Name, "Unit.Name should be set")
	assert.True(t, len(unit.Tags) >= 0, "Unit.Tags should be set")

	assert.NotNil(t, unit.Drawable.Init, "Drawable.Init should be set")
	assert.NotNil(t, unit.Drawable.Wipe, "Drawable.Wipe should be set")
	assert.NotNil(t, unit.Drawable, "Drawable.Draw should be set")
}
