package adapter

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render"
	"github.com/Rafael24595/go-terminal/engine/render/marker"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func JoinLines(inner render.RawAdapter) render.Adapter {
	return func(lines []text.Line, size winsize.Winsize) string {
		buffer := inner(lines, size)
		buffer = normalize(buffer, size.Rows)
		return strings.Join(buffer, "\n")
	}
}

func WithPadding(
	transform func(winsize.Winsize) winsize.Winsize,
	inner render.RawAdapter,
) render.Adapter {
	return func(lines []text.Line, size winsize.Winsize) string {
		r := transform(size)

		rows := min(r.Rows, size.Rows)
		cols := min(r.Cols, size.Cols)
		newSize := winsize.New(rows, cols)

		content := inner(lines, newSize)
		content = normalize(content, rows)

		diffRows := math.SubClampZero(size.Rows, rows)
		diffCols := math.SubClampZero(size.Cols, cols)

		topPadding := diffRows / 2
		leftPadding := diffCols / 2

		buffer := make([]string, size.Rows)
		for i := range size.Rows {
			buffer[i] = helper.FillRight(marker.DefaultPaddingText, int(size.Cols))
		}

		index := topPadding
		for _, line := range content {
			fixed := helper.Right(marker.DefaultPaddingText, int(leftPadding))
			buffer[index] = helper.Right(fixed+line, int(size.Cols))
			index += 1
		}

		return strings.Join(buffer, "\n")
	}
}

func normalize(lines []string, rows winsize.Rows) []string {
	buffer := make([]string, rows)
	copy(buffer, lines)
	return buffer
}
