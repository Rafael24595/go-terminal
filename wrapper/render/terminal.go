package wrapper_render

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

func TerminalRender(lines []text.Line, size terminal.Winsize) string {
	buffer := make([]string, size.Rows)

	body := terminalRenderBuffer(lines, size)

	copy(buffer, body)

	return strings.Join(buffer, "\n")
}

func terminalRenderBuffer(lines []text.Line, size terminal.Winsize) []string {
	buffer := make([]string, len(lines))

	ctx := style.LayoutContext{
		Cols: int(size.Cols),
	}

	for i, line := range lines {
		measure := text.LineFragmentsMeasurWithContext(line, ctx)
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

func renderLineFragments(line text.Line, size terminal.Winsize) string {
	var buffer strings.Builder

	fragments := ""
	atomStyles := style.AtmNone

	lineSize := terminal.Winsize{
		Rows: size.Rows,
		Cols: size.Cols,
	}

	for _, f := range line.Text {
		spec := applySpecStyles(f.Spec, lineSize, f.Text, f.Len())

		fragSize := text.FragmentMeasure(f)
		lineSize.Cols = math.SubClampZero(lineSize.Cols, uint16(fragSize))

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
