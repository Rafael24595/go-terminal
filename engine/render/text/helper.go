package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

func MaxLineMeasure(cols winsize.Cols, lines ...Line) winsize.Cols {
	size := winsize.Cols(0)
	for _, l := range lines {
		measure := FragmentMeasure(cols, l.Text...)
		size = max(size, measure)
	}
	return size
}

func FragmentsFromString(text ...string) []Fragment {
	fragments := make([]Fragment, len(text))
	for i, v := range text {
		fragments[i] = *NewFragment(v)
	}
	return fragments
}

func LineJump() *Line {
	return &Line{
		Text: FragmentsFromString(""),
		Spec: style.SpecFromKind(style.SpcKindFill),
	}
}

func ApplyLineSpec(style style.Spec, lines ...Line) []Line {
	for i := range lines {
		lines[i].SetSpec(style)
	}
	return lines
}

func LinesHasAtom(atom style.Atom, lines ...Line) bool {
	for _, line := range lines {
		if FragsHasAtom(atom, line.Text...) {
			return true
		}
	}
	return false
}

func FragsHasAtom(atom style.Atom, frags ...Fragment) bool {
	for _, v := range frags {
		if v.Atom.HasAny(atom) {
			return true
		}
	}
	return false
}

func EraseAtom(atom style.Atom, lines ...Line) bool {
	for _, line := range lines {
		for _, v := range line.Text {
			v.Atom = style.EraseAtom(v.Atom, atom)
		}
	}
	return false
}

func CloneLines(lines ...Line) []Line {
	clones := make([]Line, len(lines))
	for i, v := range lines {
		clones[i] = *v.Clone()
	}
	return clones
}
