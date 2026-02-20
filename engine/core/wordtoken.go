package core

import (
	"unicode"

	"github.com/Rafael24595/go-terminal/engine/core/style"
)

type WordToken struct {
	Text []Fragment
}

func WordTokenFromFragments(fragments ...Fragment) WordToken {
	return WordToken{
		Text: fragments,
	}
}

func (t WordToken) Size() int {
	size := 0
	for _, v := range t.Text {
		size += v.Len()
	}
	return size
}

func TokenizeLineWords(line Line) []WordToken {
	tokens := make([]WordToken, 0)
	fragments := make([]Fragment, 0)

	flush := func(pf, cf Fragment) Fragment {
		if len(cf.Text) > 0 {
			fragments = append(fragments, cf)
		}

		if len(fragments) > 0 {
			token := WordToken{
				Text: fragments,
			}

			tokens = append(tokens, token)
			fragments = make([]Fragment, 0)
		}

		return EmptyFragment().
			AddAtom(pf.Atom).
			AddSpec(pf.Spec)
	}

	inSpace := false

	for _, frag := range line.Text {
		if frag.Spec.Kind() != style.SpcKindNone {
			tokens = append(tokens, WordToken{
				Text: []Fragment{frag},
			})

			continue
		}

		fragment := EmptyFragment().
			AddAtom(frag.Atom).
			AddSpec(frag.Spec)

		for _, r := range frag.Text {
			if unicode.IsSpace(r) {
				if !inSpace {
					fragment = flush(frag, fragment)
				}

				fragment.Text += string(r)

				inSpace = true

				continue
			}

			if inSpace {
				fragment = flush(frag, fragment)
			}

			fragment.Text += string(r)

			inSpace = false
		}

		fragments = append(fragments, fragment)
	}

	if len(fragments) > 0 {
		flush(Fragment{}, Fragment{})
	}

	return tokens
}

func SplitLongToken(word WordToken, cols int, current Line, width int) (Line, []Line, int) {
	emmited := make([]Line, 0)
	if cols == 0 {
		emmited = append(emmited, LineFromFragments(word.Text...))
		return current, emmited, 0
	}

	fragments := append([]Fragment{}, word.Text...)

	flush := func() {
		emmited = append(emmited, current)
		current = LineFromSpec(current.Spec)
		width = 0
	}

	for len(fragments) > 0 {
		remaining := cols - width
		if remaining == 0 {
			flush()
			continue
		}

		fragment := fragments[0]
		size := fragment.Len()

		if size <= remaining {
			current.Text = append(current.Text, fragment)
			width += size

			fragments = fragments[1:]

			continue
		}

		taken, rest := takeFromFragment(fragment, remaining)

		current.Text = append(current.Text, taken)
		width += taken.Len()

		fragments[0] = rest

		flush()
	}

	return current, emmited, width
}

func takeFromFragment(f Fragment, n int) (Fragment, Fragment) {
	runes := []rune(f.Text)

	if n >= len(runes) {
		return f, EmptyFragment().AddAtom(f.Atom).AddSpec(f.Spec)
	}

	taken := EmptyFragment().AddAtom(f.Atom).AddSpec(f.Spec)
	taken.Text = string(runes[:n])

	rest := EmptyFragment().AddAtom(f.Atom).AddSpec(f.Spec)
	rest.Text = string(runes[n:])

	return taken, rest
}
