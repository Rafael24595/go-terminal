package wrapper_render

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

func TerminalRender(vm core.ViewModel, size terminal.Winsize) string {
	buffer := make([]string, size.Rows)

	for i, l := range vm.Lines {
		bufferLine := make([]string, 0)
		for _, f := range l.Text {
			bufferLine = append(bufferLine,
				terminalRenderFragment(f),
			)
		}

		buffer[i] = terminalRenderLine(
			l,
			size,
			bufferLine,
		)
	}

	return strings.Join(buffer, "\n")
}

func terminalRenderFragment(f core.Fragment) string {
	text := f.Text
	for _, s := range f.Styles {
		switch s {
		case core.Bold:
			text = "\033[1m" + text
		}
	}
	return text
}

func terminalRenderLine(l core.Line, size terminal.Winsize, buffer []string) string {
	line := strings.Join(buffer, " ")

	switch l.Padding.Padding {
	case core.Center:
		return helper.Center(line, int(size.Cols))
	case core.Left:
		return helper.Left(line, int(size.Cols))
	case core.Right:
		return helper.Right(line, int(size.Cols))
	case core.Fill:
		return helper.Fill(line, int(size.Cols))
	case core.Custom:
		line = helper.Left(line, int(l.Padding.Left))
		return helper.Right(line, int(l.Padding.Right))
	}

	assert.AssertfFalse(true, "undefined padding mode %d", l.Padding.Padding)

	return line
}
