package pager

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestEnginePage(t *testing.T) {
	engine := EnginePage()

	ctx := &draw.DrawContext{
		Size: winsize.Winsize{Rows: 3},
	}

	state := &draw.DrawState{
		Buffer: []text.Line{{}, {}, {}},
		Cursor: 2,
		Page:   1,
		Focus:  true,
	}

	result := engine.Func(ctx, state)

	assert.Len(t, 3, result.Buffer)
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 0, result.Cursor)
	assert.False(t, result.Focus)
}

func TestEnginePage_AlwaysResetsBuffer(t *testing.T) {
	engine := EnginePage()

	ctx := &draw.DrawContext{
		Size: winsize.Winsize{Rows: 2},
	}

	state := &draw.DrawState{
		Buffer: []text.Line{{}, {}},
	}

	engine.Func(ctx, state)
	engine.Func(ctx, state)

	assert.Equal(t, 2, state.Page)
}

func TestEngineScroll(t *testing.T) {
	engine := EngineScroll()

	ctx := &draw.DrawContext{
		Size: winsize.Winsize{Rows: 3},
	}

	state := &draw.DrawState{
		Buffer: []text.Line{
			*text.NewLine("A"),
			*text.NewLine("B"),
			*text.NewLine("C"),
		},
		Cursor: 2,
		Page:   1,
		Focus:  true,
	}

	result := engine.Func(ctx, state)

	assert.Equal(t, "B", text.LineToString(&result.Buffer[0]))
	assert.Equal(t, "C", text.LineToString(&result.Buffer[1]))
	assert.Equal(t, "", text.LineToString(&result.Buffer[2]))
	assert.Equal(t, 1, result.Cursor)
	assert.False(t, result.Focus)
}

func TestEngineScroll_CursorNeverNegative(t *testing.T) {
	engine := EngineScroll()

	state := &draw.DrawState{
		Cursor: 0,
	}

	result := engine.Func(&draw.DrawContext{}, state)

	assert.Equal(t, 0, result.Cursor)
}
