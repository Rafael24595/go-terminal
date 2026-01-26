package wrapper_render

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/render"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type FixedRender struct {
	render  render.Render
	maxRows uint16
	maxCols uint16
}

func NewFixed(render render.Render, maxRows, maxCols uint16) FixedRender {
	return FixedRender{
		render:  render,
		maxRows: maxRows,
		maxCols: maxCols,
	}
}

func (r FixedRender) ToRender() render.Render {
	return render.Render{
		Render: r.Render,
	}
}

func (r FixedRender) Render(lines []core.Line, size terminal.Winsize) string {
	rows := min(r.maxRows, size.Rows)
	cols := min(r.maxCols, size.Cols)
	newSize := terminal.NewWinsize(rows, cols)

	content := r.render.Render(lines, newSize)

	diffRows := int(size.Rows - rows)
	diffCols := int(size.Cols - cols)

	renderedLines := strings.Split(content, "\n")

	topPadding := diffRows / 2
	leftPadding := diffCols / 2

	buffer := make([]string, size.Rows)

	for i := range size.Rows {
		buffer[i] = helper.Fill(" ", int(size.Cols))
	}

	index := topPadding
	for _, line := range renderedLines {
		fixed := helper.Right(" ", leftPadding)
		buffer[index] = helper.Right(fixed + line, int(size.Cols))
		index += 1
	}

	return strings.Join(buffer, "\n")
}
