package inputline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestInputLine_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).ToUnit()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestNewInputLine_DefaultPrompt(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	input := New(mock.ToUnit())

	assert.Equal(t, input.prompt, marker.DefaultPromptText)
}

func TestNewInputLine_NoContent_ReturnsPromptOnly(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Status: false,
		Lines:  make([]text.Line, 0),
	}

	unit := New(mock.ToUnit()).ToUnit()

	unit.Drawable.Init()
	lines, status := unit.Drawable.Draw(winsize.Winsize{
		Rows: 5,
	})

	assert.False(t, status)
	assert.Len(t, 1, lines)
	assert.Equal(t, marker.DefaultPromptText, text.LineToString(&lines[0]))
}

func TestNewInputLine_WithSingleLine_AddsPrompt(t *testing.T) {
	frag := text.FragmentsFromString("golang")

	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag...),
		},
	}

	unit := New(mock.ToUnit()).ToUnit()

	unit.Drawable.Init()
	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Cols: 10,
		Rows: 5,
	})

	assert.Len(t, 1, lines)
	assert.Equal(t, marker.DefaultPromptText+" golang", text.LineToString(&lines[0]))
}

func TestNewInputLine_MultipleDrawCalls_AccumulatesLines(t *testing.T) {
	frag1 := text.FragmentsFromString("ziglang")
	frag2 := text.FragmentsFromString("golang")

	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag1...),
			*text.LineFromFragments(frag2...),
		},
	}

	unit := New(mock.ToUnit()).ToUnit()

	unit.Drawable.Init()
	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Cols: 10,
		Rows: 5,
	})

	assert.Len(t, 2, lines)

	assert.Equal(t, marker.DefaultPromptText+" ziglang", text.LineToString(&lines[0]))
	assert.Equal(t, "golang", text.LineToString(&lines[1]))
}
