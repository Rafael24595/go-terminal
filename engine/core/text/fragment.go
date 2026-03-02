package text

import (
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core/style"
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

func EmptyFragment() Fragment {
	return NewFragment("")
}

func EmptyFragmentFrom(frag Fragment) Fragment {
	return FragmentFrom("", frag)
}

func (f Fragment) AddAtom(styles ...style.Atom) Fragment {
	f.Atom = style.MergeAtom(styles...)
	return f
}

func (f Fragment) AddSpec(styles ...style.Spec) Fragment {
	f.Spec = style.MergeSpec(styles...)
	return f
}

func (f Fragment) Len() int {
	return utf8.RuneCountInString(f.Text)
}

func FragmentMeasure(frag Fragment) int {
	return style.SpecMeasure(frag.Spec, frag.Len())
}

func FragmentMeasureWithContext(frag Fragment, ctx style.LayoutContext) int {
	return style.SpecMeasureWithContext(frag.Spec, frag.Len(), ctx)
}
