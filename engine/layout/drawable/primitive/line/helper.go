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

func NextIndexedWrappedLine(cols winsize.Cols, lines []text.Line, meta indexMeta) (*text.Line, []text.Line) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]text.Line, 0)
	}

	var prefix string
	if lines[0].Order != 0 {
		order := int(lines[0].Order)
		prefix = meta.header(order)
		lines[0].Order = 0
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

func computeIndexMeta(lines []text.Line) *indexMeta {
	size := winsize.Cols(0)

	for _, line := range lines {
		if line.Order == 0 {
			continue
		}

		digits := math.Digits(line.Order)
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
