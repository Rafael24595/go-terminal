package pipeline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func mockInitStep(s winsize.Winsize, d drawable.Drawable) ([]text.Line, bool) {
	return d.Draw(s)
}

func mockDataStep(_ winsize.Winsize, l []text.Line, s bool) ([]text.Line, bool) {
	return l, s
}

func TestPipeline_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := NewPipelineDrawable(mock.ToDrawable()).
		SetInitStep(mockInitStep).
		ToDrawable()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestPipeline_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewPipelineDrawable(mock.ToDrawable()).
		SetInitStep(mockInitStep)

	dw := bd.ToDrawable()
	dw.Init()

	assert.Panic(t, func() {
		bd.SetInitStep(mockInitStep)
	})

	assert.Panic(t, func() {
		bd.PushDataSteps(mockDataStep)
	})
}

func TestPipeline_ReturnBaseIfNils(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := NewPipelineDrawable(mock.ToDrawable()).
		ToDrawable()

	assert.Equal(t, drawable_test.NameMockDrawable, dw.Name)

	dw = NewPipelineDrawable(mock.ToDrawable()).
		SetInitStep(mockInitStep).
		ToDrawable()

	assert.Equal(t, Name, dw.Name)
}

func TestPipeline_InitStepTransformation(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.NewLine("base_02"),
			*text.NewLine("base_03"),
		},
		Status: true,
	}

	mockLine := text.NewLine("mock_line_01")
	bd := NewPipelineDrawable(mock.ToDrawable()).
		SetInitStep(func(_ winsize.Winsize, _ drawable.Drawable) ([]text.Line, bool) {
			return []text.Line{*mockLine}, false
		})

	dw := bd.ToDrawable()
	dw.Init()

	lines, status := dw.Draw(winsize.Winsize{})

	assert.Len(t, 1, lines)
	assert.False(t, status)
	assert.Equal(t, text.LineToString(mockLine), text.LineToString(&lines[0]))
}

func TestPipeline_DataStepsChain(t *testing.T) {
	baseLine := text.NewLine("base_01")
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*baseLine,
		},
		Status: true,
	}

	mockLine1 := text.NewLine("mock_line_01")
	mockLine2 := text.NewLine("mock_line_02")
	bd := NewPipelineDrawable(mock.ToDrawable()).
		PushDataSteps(
			func(_ winsize.Winsize, l []text.Line, s bool) ([]text.Line, bool) {
				return append(l, *mockLine1), s
			},
			func(_ winsize.Winsize, l []text.Line, s bool) ([]text.Line, bool) {
				return append(l, *mockLine2), !s
			},
		)

	dw := bd.ToDrawable()
	dw.Init()

	lines, status := dw.Draw(winsize.Winsize{})

	assert.Len(t, 3, lines)
	assert.False(t, status)

	assert.Equal(t, text.LineToString(baseLine), text.LineToString(&lines[0]))
	assert.Equal(t, text.LineToString(mockLine1), text.LineToString(&lines[1]))
	assert.Equal(t, text.LineToString(mockLine2), text.LineToString(&lines[2]))
}
