package wrapper_render

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TerminalRawRender(lines []text.Line, size winsize.Winsize) []string {
	buffer := make([]string, len(lines))

	for i, line := range lines {
		measure := text.FragmentMeasure(int(size.Cols), line.Text...)
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

func renderLineFragments(line text.Line, size winsize.Winsize) string {
	var buffer strings.Builder

	fragments := ""
	atomStyles := style.AtmNone

	lineSize := winsize.New(
		size.Rows,
		size.Cols,
	)

	for _, f := range line.Text {
		spec := applySpecStyles(f.Spec, lineSize, f.Text, f.Size())

		fragSize := text.FragmentMeasure(int(size.Cols), f)
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
