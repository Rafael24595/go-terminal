package inline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestInline_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := InlineDrawableFromDrawables(mock.ToDrawable())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestInline_LazyInit(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewInlineDrawable(mock.ToDrawable())

	assert.False(t, bd.lazyLoaded)
	assert.True(t, 
		drawable_test.Helper_IsDrawableNil(t, bd.drawable),
	)

	bd.init()
	bd.draw(terminal.Winsize{})

	assert.True(t, bd.lazyLoaded)
	drawable_test.Helper_ToDrawable(t, bd.drawable)
}

func TestInline_JoinsChildren(t *testing.T) {
	mock1 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("go"),
		},
	}

	mock2 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("lang"),
		},
	}

	d := NewInlineDrawable(
		mock1.ToDrawable(),
		mock2.ToDrawable(),
	)

	dr := d.ToDrawable()

	dr.Init()

	lines, _ := dr.Draw(terminal.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.Len(t, 1, lines)
	assert.Equal(t, "golang", text.LineToString(&lines[0]))
}

func TestInline_JoinsChildrenWithSeparator(t *testing.T) {
	mock1 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("golang"),
		},
	}

	mock2 := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("ziglang"),
		},
	}

	d := NewInlineDrawable(
		mock1.ToDrawable(),
		mock2.ToDrawable(),
	)

	d.Separator(" | ")

	dr := d.ToDrawable()

	dr.Init()

	lines, _ := dr.Draw(terminal.Winsize{
		Rows: 3,
		Cols: 16,
	})

	assert.Len(t, 1, lines)
	assert.Equal(t, "golang | ziglang", text.LineToString(&lines[0]))
}

func TestInline_MultipleLines(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("go"),
			*text.NewLine("lang"),
		},
	}

	d := NewInlineDrawable(
		mock.ToDrawable(),
	)

	dr := d.ToDrawable()

	d.Separator(" | ")

	dr.Init()

	lines, _ := dr.Draw(terminal.Winsize{
		Rows: 3,
		Cols: 9,
	})

	assert.Len(t, 1, lines)
	assert.Equal(t, "go | lang", text.LineToString(&lines[0]))
}

func TestInline_Empty(t *testing.T) {
	d := NewInlineDrawable()

	dr := d.ToDrawable()

	dr.Init()

	lines, _ := dr.Draw(terminal.Winsize{})

	assert.Len(t, 0, lines)
}
