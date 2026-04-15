package mapdrw

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func mockDrawInputMap(w terminal.Winsize) terminal.Winsize {
	return w
}

func TestMap_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := NewMapDrawable(mock.ToDrawable()).
		SetDrawInputMap(mockDrawInputMap).
		ToDrawable()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}

func TestMap_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	bd := NewMapDrawable(mock.ToDrawable()).
		SetDrawInputMap(mockDrawInputMap)
		
	dw := bd.ToDrawable()
	dw.Init()

	assert.Panic(t, func() {
		bd.SetDrawInputMap(mockDrawInputMap)
	})
}

func TestMap_ReturnBaseIfNils(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := NewMapDrawable(mock.ToDrawable()).
		ToDrawable()

	assert.Equal(t, drawable_test.NameMockDrawable, dw.Name)

	dw = NewMapDrawable(mock.ToDrawable()).
		SetDrawInputMap(mockDrawInputMap).
		ToDrawable()

	assert.Equal(t, NameMapDrawable, dw.Name)
}

func TestMap_InputMapping(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	mockSize := terminal.Winsize{
		Rows: 50,
		Cols: 50,
	}

	dw := NewMapDrawable(mock.ToDrawable()).
		SetDrawInputMap(PredFixedWinsize(mockSize)).
		ToDrawable()

	size := terminal.Winsize{
		Rows: 25,
		Cols: 25,
	}

	dw.Init()
	dw.Draw(size)

	assert.Equal(t, mockSize.Rows, mock.Size.Rows)
	assert.Equal(t, mockSize.Cols, mock.Size.Cols)
}

func TestMap_OutputMapping(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("test"),
		},
	}

	dw := NewMapDrawable(mock.ToDrawable()).
		SetDrawOutputMap(func(w terminal.Winsize, d drawable.Drawable) ([]text.Line, bool) {
			lines, _ := d.Draw(w)
			return []text.Line{
				*lines[0].UnshiftFragments(
					*text.NewFragment("mocked"),
					*text.NewFragment(" "),
				),
			}, true
		}).
		ToDrawable()

	size := terminal.Winsize{
		Rows: 25,
		Cols: 25,
	}

	dw.Init()
	lines, _ := dw.Draw(size)

	assert.Equal(t, "mocked test", text.LineToString(&lines[0]))
}
