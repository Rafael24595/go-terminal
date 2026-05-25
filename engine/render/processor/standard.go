package processor

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Standard struct {
	atom styler.Atom
	spec styler.Spec
}

func New(atom styler.Atom, spec styler.Spec) Standard {
	return Standard{
		atom: atom,
		spec: spec,
	}
}

func (r Standard) Render(lines []text.Line, size winsize.Winsize) []string {
	buffer := make([]string, len(lines))

	for i, line := range lines {
		measure := text.FragmentMeasure(size.Cols, line.Text...)
		styled := r.renderLineFragments(line, size)

		buffer[i] = r.spec.Apply(
			line.Spec,
			size,
			styled,
			measure,
		)
	}

	return buffer
}

func (r Standard) renderLineFragments(line text.Line, size winsize.Winsize) string {
	var buffer strings.Builder

	fragments := ""
	atomStyles := style.AtmNone

	lineSize := winsize.New(
		size.Rows,
		size.Cols,
	)

	for _, f := range line.Text {
		spec := r.spec.Apply(f.Spec, lineSize, f.Text, f.Size())

		fragSize := text.FragmentMeasure(size.Cols, f)
		lineSize.Cols = lineSize.Cols.Sub(fragSize)

		if atomStyles != f.Atom && len(fragments) != 0 {
			atom := r.atom.Apply(fragments, atomStyles)
			buffer.WriteString(atom)

			fragments = spec
			atomStyles = f.Atom

			continue
		}

		fragments += spec
		atomStyles = f.Atom
	}

	if len(fragments) != 0 {
		atom := r.atom.Apply(fragments, atomStyles)
		buffer.WriteString(atom)
	}

	return buffer.String()
}

