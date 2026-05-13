package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

func TestWrapNextLine_FitWithMeta(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: "    ",
	}

	layout := wrap.NormalizeLines(
		*text.NewLine("golang").SetOrder(1),
	)

	got, remain := NextIndexedLine(10, layout, meta)

	assert.Equal(t, "1 | golang", text.LineToString(got))

	assert.Len(t, 0, remain)
}

func TestWrapNextLine_SplitWithMeta(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	layout := wrap.NormalizeLines(
		*text.NewLine("golang rust").SetOrder(1),
	)

	got, remain := NextIndexedLine(10, layout, meta)

	assert.Equal(t, "1 | golang", text.LineToString(got))
	assert.Len(t, 1, remain)

	got, remain = NextIndexedLine(10, remain, meta)

	assert.Equal(t, "  |  rust", text.LineToString(got))
	assert.Len(t, 0, remain)
}

func TestWrapNextLine_IndexShouldBeLesser(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	layout := wrap.NewLayoutLine(
		text.NewLine("golang").SetOrder(1),
	)

	assert.Panic(t, func() {
		NextIndexedLine(4, []wrap.LayoutLine{*layout}, meta)
	})
}
