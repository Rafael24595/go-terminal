package wrapper_render

import (
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

func terminalRenderBuffer(lines []core.Line, size terminal.Winsize) []string {
	buffer := make([]string, len(lines))

	for i, l := range lines {
		styled := renderLineFragments(l)

		buffer[i] = terminalRenderLine(
			lines,
			i,
			size,
			styled,
		)
	}

	return buffer
}

func renderLineFragments(l core.Line) string {
	var line strings.Builder

	fragments := ""
	styles := core.None

	for _, f := range l.Text {
		if styles != f.Styles && len(fragments) != 0 {
			line.WriteString(applystyles(fragments, styles))

			fragments = ""

			fragments += f.Text
			styles = f.Styles
			continue
		}

		fragments += f.Text
		styles = f.Styles
	}

	if len(fragments) != 0 {
		line.WriteString(applystyles(fragments, styles))
	}

	return line.String()
}

func applystyles(text string, styles ...core.Style) string {
	merged := core.MergeStyles(styles...)

	if merged.HasAny(core.Lower) {
		text = strings.ToLower(text)
	}

	if merged.HasAny(core.Upper) {
		text = strings.ToUpper(text)
	}

	if merged.HasAny(core.Bold) {
		text = Bold + text + NoBold
	}

	if merged.HasAny(core.Select) {
		text = Reverse + text + NoReverse
	}

	return text
}

func terminalRenderLine(lines []core.Line, index int, size terminal.Winsize, line string) string {
	padd := lines[index].Padding

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

	assert.AssertFalse(true, "undefined padding mode %d", padd.Padding)

	return line
}
