package text

import "github.com/Rafael24595/go-terminal/engine/render/style"

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

func HasFocus(line *Line) bool {
	for _, v := range line.Text {
		if v.Atom.HasAny(style.AtmFocus) {
			return true
		}
	}
	return false
}
