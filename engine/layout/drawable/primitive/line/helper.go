package line

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const separator = " | "

func NextIndexedLine(cols winsize.Cols, lines []wrap.LayoutLine, meta indexMeta) (*text.Line, []wrap.LayoutLine) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]wrap.LayoutLine, 0)
	}

	var prefix string
	if lines[0].Source.Order != 0 {
		order := int(lines[0].Source.Order)
		prefix = meta.header(order)
		lines[0].Source.Order = 0
	} else {
		prefix = meta.body()
	}

	fixedCols := cols.Clamp(meta.totalWidth)

	assert.True(fixedCols > 0, "index prefix should be lesser than line size")

	cursor, rest := wrap.NextLine(fixedCols, lines)
	if cursor != nil {
		cursor.UnshiftFragments(
			*text.NewFragment(prefix),
		)
	}

	return cursor, rest
}

func computeIndexMeta(lines []wrap.LayoutLine) *indexMeta {
	size := winsize.Cols(0)

	for _, line := range lines {
		if line.Source.Order == 0 {
			continue
		}

		digits := math.Digits(line.Source.Order)
		size = max(size, winsize.Cols(digits))
	}

	if size == 0 {
		return nil
	}

	return &indexMeta{
		sufix:      separator,
		prefixBody: helper.FillRight(marker.DefaultPaddingText, size),
		digits:     uint16(size),
		totalWidth: size + runes.Measure(separator),
	}
}
