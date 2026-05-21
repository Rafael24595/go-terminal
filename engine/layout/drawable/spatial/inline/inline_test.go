package inline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestInline_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := UnitFromUnits(mock.ToUnit())
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestInline_JoinsChildren(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("go"),
		},
	}
	mock2 := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("lang"),
		},
	}

	unit := New(
		mock1.ToUnit(),
		mock2.ToUnit(),
	).ToUnit()

	unit.Drawable.Init()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.Len(t, 1, lines)
	assert.Equal(t, "golang", text.LineToString(&lines[0]))
}

func TestInline_JoinsChildrenWithSeparator(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("golang"),
		},
	}
	mock2 := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("ziglang"),
		},
	}

	unit := New(
		mock1.ToUnit(),
		mock2.ToUnit(),
	).Separator(" | ").ToUnit()

	unit.Drawable.Init()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 16,
	})

	assert.Len(t, 1, lines)
	assert.Equal(t, "golang | ziglang", text.LineToString(&lines[0]))
}

func TestInline_MultipleLines(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("go"),
			*text.NewLine("lang"),
		},
	}

	unit := New(
		mock.ToUnit(),
	).Separator(" | ").ToUnit()

	unit.Drawable.Init()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 9,
	})

	assert.Len(t, 1, lines)
	assert.Equal(t, "go | lang", text.LineToString(&lines[0]))
}

func TestInline_Empty(t *testing.T) {
	unit := New().ToUnit()

	unit.Drawable.Init()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{})

	assert.Len(t, 0, lines)
}
