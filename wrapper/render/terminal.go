package wrapper_render

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/style"
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

	for i, line := range lines {
		styled := renderLineFragments(line, size)

		buffer[i] = renderLine(
			lines,
			i,
			size,
			styled,
		)
	}

	return buffer
}

func renderLineFragments(line core.Line, size terminal.Winsize) string {
	var buffer strings.Builder

	fragments := ""
	styles := style.AtmNone

	flush := func(f core.Fragment) {
		atom := applyAtomStyles(fragments, styles)
		spec, _ := applyVariantStyles(f.Spec, size, atom)

		buffer.WriteString(spec)

		fragments = ""

		fragments += f.Text
		styles = f.Atom
	}

	for _, f := range line.Text {
		if styles != f.Atom && len(fragments) != 0 {
			flush(f)
			continue
		}

		fragments += f.Text
		styles = f.Atom
	}

	if len(fragments) != 0 {
		flush(core.EmptyFragment())
	}

	return buffer.String()
}

func renderLine(lines []core.Line, index int, size terminal.Winsize, line string) string {
	if line, ok := applyLineVariantStyles(lines, index, size, line); ok {
		return line
	}

	styl := lines[index].Spec
	line, _ = applyVariantStyles(styl, size, line)

	return line
}
