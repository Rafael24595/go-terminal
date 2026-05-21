package pipeline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func mockInitStep(s winsize.Winsize, d drawable.Unit) drawable.Unit {
	return d
}

func mockDrawStep(s winsize.Winsize, d drawable.Unit) ([]text.Line, bool) {
	return d.Drawable.Draw(s)
}

func mockDataStep(_ winsize.Winsize, _ drawable.Unit, l []text.Line, s bool) ([]text.Line, bool) {
	return l, s
}

func TestPipeline_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).
		SetDrawStep(mockDrawStep).
		ToUnit()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestPipeline_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).
		SetDrawStep(mockDrawStep)

	unit.ToUnit().Drawable.Init()

	assert.Panic(t, func() {
		unit.PushInitSteps(mockInitStep)
	})

	assert.Panic(t, func() {
		unit.SetDrawStep(mockDrawStep)
	})

	assert.Panic(t, func() {
		unit.PushDataSteps(mockDataStep)
	})
}

func TestPipeline_ReturnBaseIfNils(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).
		ToUnit()

	assert.Equal(t, drawable_test.NameMockUnit, unit.Name)

	unit = New(mock.ToUnit()).
		SetDrawStep(mockDrawStep).
		ToUnit()

	assert.Equal(t, Name, unit.Name)
}

func TestPipeline_InitStepTransformation(t *testing.T) {
	mock1 := &drawable_test.MockUnit{}

	mock2 := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.NewLine("base_02"),
		},
		Status: true,
	}

	unit := New(mock1.ToUnit()).
		PushInitSteps(func(_ winsize.Winsize, _ drawable.Unit) drawable.Unit {
			return mock2.ToUnit()
		}).
		ToUnit()

	unit.Drawable.Init()

	lines, status := unit.Drawable.Draw(winsize.Winsize{})

	assert.Len(t, 2, lines)
	assert.True(t, status)
	assert.Equal(t, text.LineToString(&mock2.Lines[0]), text.LineToString(&lines[0]))
}

func TestPipeline_DrawStepTransformation(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.NewLine("base_02"),
			*text.NewLine("base_03"),
		},
		Status: true,
	}

	mockLine := text.NewLine("mock_line_01")
	unit := New(mock.ToUnit()).
		SetDrawStep(func(_ winsize.Winsize, _ drawable.Unit) ([]text.Line, bool) {
			return []text.Line{*mockLine}, false
		}).
		ToUnit()

	unit.Drawable.Init()

	lines, status := unit.Drawable.Draw(winsize.Winsize{})

	assert.Len(t, 1, lines)
	assert.False(t, status)
	assert.Equal(t, text.LineToString(mockLine), text.LineToString(&lines[0]))
}

func TestPipeline_DataStepsChain(t *testing.T) {
	baseLine := text.NewLine("base_01")
	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*baseLine,
		},
		Status: true,
	}

	mockLine1 := text.NewLine("mock_line_01")
	mockLine2 := text.NewLine("mock_line_02")
	
	unit := New(mock.ToUnit()).
		PushDataSteps(
			func(_ winsize.Winsize, _ drawable.Unit, l []text.Line, s bool) ([]text.Line, bool) {
				return append(l, *mockLine1), s
			},
			func(_ winsize.Winsize, _ drawable.Unit, l []text.Line, s bool) ([]text.Line, bool) {
				return append(l, *mockLine2), !s
			},
		).ToUnit()

	unit.Drawable.Init()

	lines, status := unit.Drawable.Draw(winsize.Winsize{})

	assert.Len(t, 3, lines)
	assert.False(t, status)

	assert.Equal(t, text.LineToString(baseLine), text.LineToString(&lines[0]))
	assert.Equal(t, text.LineToString(mockLine1), text.LineToString(&lines[1]))
	assert.Equal(t, text.LineToString(mockLine2), text.LineToString(&lines[2]))
}
