package inline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestInline_ToDrawable(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := InlineDrawableFromDrawables(mock.ToDrawable())
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestInlineDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewInlineDrawable(mock.ToDrawable())

	assert.Panic(t, func() {
		bd.draw()
	})
}

func TestInlineDrawable_JoinsChildren(t *testing.T) {
	mock1 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			text.LineFromString("go"),
		},
	}

	mock2 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			text.LineFromString("lang"),
		},
	}

	d := NewInlineDrawable(
		mock1.ToDrawable(),
		mock2.ToDrawable(),
	)

	dr := d.ToDrawable()

	dr.Init(terminal.Winsize{})

	lines, _ := dr.Draw()

	assert.Len(t, 1, lines)
	assert.Equal(t, "golang", text.LineToString(lines[0]))
}

func TestInlineDrawable_JoinsChildrenWithSeparator(t *testing.T) {
	mock1 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			text.LineFromString("golang"),
		},
	}

	mock2 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			text.LineFromString("ziglang"),
		},
	}

	d := NewInlineDrawable(
		mock1.ToDrawable(),
		mock2.ToDrawable(),
	)

	d.Separator(" | ")

	dr := d.ToDrawable()

	dr.Init(terminal.Winsize{
		Cols: 16,
	})

	lines, _ := dr.Draw()

	assert.Len(t, 1, lines)
	assert.Equal(t, "golang | ziglang", text.LineToString(lines[0]))
}

func TestInlineDrawable_MultipleLines(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			text.LineFromString("go"),
			text.LineFromString("lang"),
		},
	}

	d := NewInlineDrawable(
		mock.ToDrawable(),
	)

	dr := d.ToDrawable()

	d.Separator(" | ")

	dr.Init(terminal.Winsize{
		Cols: 9,
	})

	lines, _ := dr.Draw()

	assert.Len(t, 1, lines)
	assert.Equal(t, "go | lang", text.LineToString(lines[0]))
}

func TestInlineDrawable_Empty(t *testing.T) {
	d := NewInlineDrawable()

	dr := d.ToDrawable()

	dr.Init(terminal.Winsize{})

	lines, _ := dr.Draw()

	assert.Len(t, 0, lines)
}
