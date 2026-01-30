package wrapper_render

import (
	"slices"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

func TerminalRender(lines []core.Line, size terminal.Winsize) string {
	buffer := make([]string, size.Rows)

	body := terminalRenderBuffer(lines, size)

	copy(buffer, body)

	return strings.Join(buffer, "\n")
}

func terminalRenderBuffer(ls []core.Line, size terminal.Winsize) []string {
	buffer := make([]string, len(ls))

	for i, l := range ls {
		bufferLine := make([]string, 0)
		for _, f := range l.Text {
			fragment := terminalRenderFragment(f)
			isJoin := slices.Contains(f.Styles, core.Join)

			if isJoin && len(bufferLine) > 0 {
				bufferLine[len(bufferLine)-1] += fragment
				continue
			}

			bufferLine = append(bufferLine, fragment)
		}

		buffer[i] = terminalRenderLine(
			ls,
			i,
			size,
			bufferLine,
		)
	}

	return buffer
}

func terminalRenderFragment(f core.Fragment) string {
	text := f.Text
	for _, s := range f.Styles {
		switch s {
		case core.Bold:
			text = "\033[1m" + text
		case core.Select:
			text = "\x1b[7m" + text + "\x1b[27m" 
		}
	}
	return text
}

func terminalRenderLine(lines []core.Line, index int, size terminal.Winsize, buffer []string) string {
	padd := lines[index].Padding
	line := strings.Join(buffer, " ")

	switch padd.Padding {
	case core.Center:
		return helper.Center(line, int(size.Cols))
	case core.Left:
		return helper.Left(line, int(size.Cols))
	case core.Right:
		return helper.Right(line, int(size.Cols))
	case core.Fill:
		return helper.Fill(line, int(size.Cols))
	case core.FillUp:
		cursor := index - 1
		if cursor >= len(lines) {
			return line
		}
		return helper.Fill(line, lines[cursor].Len())
	case core.FillDown:
		cursor := index + 1
		if cursor < 0 {
			return line
		}
		return helper.Fill(line, lines[cursor].Len())
	case core.Custom:
		line = helper.RepeatLeft(line, int(padd.Left))
		return helper.RepeatRight(line, int(padd.Right))
	case core.Unstyled:
		return line
	}

	assert.AssertfFalse(true, "undefined padding mode %d", padd.Padding)

	return line
}
