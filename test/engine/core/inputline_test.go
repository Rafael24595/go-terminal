package core_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTable_ToDrawable(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	input := core.NewInputLine(mock.ToDrawable())
	dw := input.ToDrawable()

	drawable_test.Helper_ToDrawable(t, dw)
}

func TestTableDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	input := core.NewInputLine(mock.ToDrawable())
	dw := input.ToDrawable()

	assert.Panic(t, func() {
		dw.Draw()
	})
}

func TestNewInputLine_DefaultPrompt(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	input := core.NewInputLine(mock.ToDrawable())

	assert.Equal(t, input.Prompt, core.DefaultPrompt)
}

func TestDraw_NoContent_ReturnsPromptOnly(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines:  core.NewLines(),
	}

	input := core.NewInputLine(mock.ToDrawable())
	drawable := input.ToDrawable()

	drawable.Init(terminal.Winsize{})
	lines, status := drawable.Draw()

	assert.False(t, status)
	assert.Len(t, 1, lines)
	assert.Equal(t, core.DefaultPrompt, core.LineToString(lines[0]))
}

func TestDraw_WithSingleLine_AddsPrompt(t *testing.T) {
	frag := core.FragmentsFromString("golang")

	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []core.Line{
			core.LineFromFragments(frag...),
		},
	}

	input := core.NewInputLine(mock.ToDrawable())
	drawable := input.ToDrawable()

	drawable.Init(terminal.Winsize{})
	lines, _ := drawable.Draw()

	assert.Len(t, 2, lines)
	assert.Equal(t, core.DefaultPrompt+" golang", core.LineToString(lines[1]))
}

func TestDraw_MultipleDrawCalls_AccumulatesLines(t *testing.T) {
	frag1 := core.FragmentsFromString("ziglang")
	frag2 := core.FragmentsFromString("golang")

	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []core.Line{
			core.LineFromFragments(frag1...),
			core.LineFromFragments(frag2...),
		},
	}

	input := core.NewInputLine(mock.ToDrawable())
	drawable := input.ToDrawable()

	drawable.Init(terminal.Winsize{})
	lines, _ := drawable.Draw()

	assert.Len(t, 3, lines)

	assert.Equal(t, core.DefaultPrompt+" ziglang", core.LineToString(lines[1]))
	assert.Equal(t, "golang", core.LineToString(lines[2]))
}
