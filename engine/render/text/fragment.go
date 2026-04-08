package text

import (
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/engine/render/style"
)

type Fragment struct {
	Text string
	Atom style.Atom
	Spec style.Spec
}

func FragmentsFromString(text ...string) []Fragment {
	fragments := make([]Fragment, len(text))
	for i, v := range text {
		fragments[i] = NewFragment(v)
	}
	return fragments
}

func NewFragment(text string) Fragment {
	return Fragment{
		Text: text,
		Atom: style.AtmNone,
		Spec: style.SpecEmpty(),
	}
}

func FragmentFrom(text string, frag Fragment) Fragment {
	return NewFragment(text).
		AddAtom(frag.Atom).
		AddSpec(frag.Spec)
}

func FragmentFromRunes(runes []rune) Fragment {
	return NewFragment(string(runes))
}

func EmptyFragment() Fragment {
	return NewFragment("")
}

func EmptyFragmentFrom(frag Fragment) Fragment {
	return FragmentFrom("", frag)
}

func (f Fragment) AddAtom(styles ...style.Atom) Fragment {
	newAtom := style.MergeAtom(styles...)
	f.Atom = style.MergeAtom(f.Atom, newAtom)
	return f
}

func (f Fragment) CutAtom(styles style.Atom) Fragment {
	f.Atom = style.EraseAtom(f.Atom, styles)
	return f
}

func (f Fragment) AddSpec(styles ...style.Spec) Fragment {
	newSpec := style.MergeSpec(styles...)
	f.Spec = style.MergeSpec(f.Spec, newSpec)
	return f
}

func (f Fragment) CutSpec(styles style.SpecKind) Fragment {
	f.Spec, _ = style.EraseSpec(f.Spec, styles)
	return f
}

func (f Fragment) Len() int {
	return runes.Measure(f.Text)
}

func FragmentMeasure(frag Fragment) int {
	return style.SpecMeasure(frag.Spec, frag.Len())
}

func FragmentMeasureWithContext(frag Fragment, ctx style.LayoutContext) int {
	return style.SpecMeasureWithContext(frag.Spec, frag.Len(), ctx)
}
