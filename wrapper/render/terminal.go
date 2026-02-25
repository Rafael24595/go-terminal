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

	ctx := style.LayoutContext{
		Cols: int(size.Cols),
	}

	for i, line := range lines {
		measure := core.LineFragmentsMeasurWithContext(line, ctx)
		styled := renderLineFragments(line, size)

		buffer[i] = applySpecStyles(
			line.Spec,
			size,
			styled,
			measure,
		)
	}

	return buffer
}

func renderLineFragments(line core.Line, size terminal.Winsize) string {
	var buffer strings.Builder

	fragments := ""
	atomStyles := style.AtmNone

	for _, f := range line.Text {
		spec := applySpecStyles(f.Spec, size, f.Text, f.Len())

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
