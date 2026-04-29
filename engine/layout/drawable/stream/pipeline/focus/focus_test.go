package focus

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestFocusInitTransformer_FocusAtStart(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.LineFromFragments(
				*text.NewFragment("base").AddAtom(style.AtmFocus),
				*text.NewFragment("_"),
				*text.NewFragment("02"),
			),
			*text.NewLine("base_03"),
		},
		Status: true,
	}

	transformer := FocusInitTransformer(
		pager.EnginePage(),
	)

	lines, status := transformer(winsize.Winsize{
		Rows: 2,
		Cols: 10,
	}, mock.ToDrawable())

	assert.Len(t, 2, lines)
	
	assert.False(t, status)
	assert.Equal(t, "base_01", text.LineToString(&lines[0]))
	assert.Equal(t, "base_02", text.LineToString(&lines[1]))
}

func TestFocusInitTransformer_FocusAtEnd(t *testing.T) {
	mock := &drawable_test.MockDrawable{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.NewLine("base_02"),
			*text.LineFromFragments(
				*text.NewFragment("base"),
				*text.NewFragment("_").AddAtom(style.AtmFocus),
				*text.NewFragment("03"),
			),
		},
		Chunk: 1,
	}

	transformer := FocusInitTransformer(
		pager.EngineScroll(),
	)

	lines, status := transformer(winsize.Winsize{
		Rows: 2,
		Cols: 10,
	}, mock.ToDrawable())

	assert.Len(t, 2, lines)
	
	assert.False(t, status)
	assert.Equal(t, "base_02", text.LineToString(&lines[0]))
	assert.Equal(t, "base_03", text.LineToString(&lines[1]))
}
