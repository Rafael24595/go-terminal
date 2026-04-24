package viewmodel_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestInputLine_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	input := viewmodel.NewInputLine(mock.ToDrawable())
	dw := input.ToDrawable()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestNewInputLine_DefaultPrompt(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	input := viewmodel.NewInputLine(mock.ToDrawable())

	assert.Equal(t, input.Prompt, marker.DefaultInputLinePrompt)
}

func TestNewInputLine_NoContent_ReturnsPromptOnly(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines:  make([]text.Line, 0),
	}

	input := viewmodel.NewInputLine(mock.ToDrawable())
	drawable := input.ToDrawable()

	drawable.Init()
	lines, status := drawable.Draw(winsize.Winsize{})

	assert.False(t, status)
	assert.Len(t, 1, lines)
	assert.Equal(t, marker.DefaultInputLinePrompt, text.LineToString(&lines[0]))
}

func TestNewInputLine_WithSingleLine_AddsPrompt(t *testing.T) {
	frag := text.FragmentsFromString("golang")

	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag...),
		},
	}

	input := viewmodel.NewInputLine(mock.ToDrawable())
	drawable := input.ToDrawable()

	drawable.Init()
	lines, _ := drawable.Draw(winsize.Winsize{})

	assert.Len(t, 2, lines)
	assert.Equal(t, marker.DefaultInputLinePrompt+" golang", text.LineToString(&lines[1]))
}

func TestNewInputLine_MultipleDrawCalls_AccumulatesLines(t *testing.T) {
	frag1 := text.FragmentsFromString("ziglang")
	frag2 := text.FragmentsFromString("golang")

	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag1...),
			*text.LineFromFragments(frag2...),
		},
	}

	input := viewmodel.NewInputLine(mock.ToDrawable())
	drawable := input.ToDrawable()

	drawable.Init()
	lines, _ := drawable.Draw(winsize.Winsize{})

	assert.Len(t, 3, lines)

	assert.Equal(t, marker.DefaultInputLinePrompt+" ziglang", text.LineToString(&lines[1]))
	assert.Equal(t, "golang", text.LineToString(&lines[2]))
}
