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
	atomStyles := style.AtmNone

	for _, f := range line.Text {
		spec, _ := applySpecStyles(f.Spec, size, f.Text)

		if atomStyles != f.Atom && len(fragments) != 0 {
			atom := applyAtomStyles(fragments, atomStyles)
			buffer.WriteString(atom)

			fragments = spec
			atomStyles = f.Atom

			continue
		}

		fragments += spec
		atomStyles = f.Atom
	}

	if len(fragments) != 0 {
		atom := applyAtomStyles(fragments, atomStyles)
		buffer.WriteString(atom)
	}

	return buffer.String()
}

func renderLine(lines []core.Line, index int, size terminal.Winsize, line string) string {
	if line, ok := applyLineSpecStyles(lines, index, size, line); ok {
		return line
	}

	styl := lines[index].Spec
	line, _ = applySpecStyles(styl, size, line)

	return line
}
